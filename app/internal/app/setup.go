package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	gRPCCategoryService "github.com/Kazzess/contracts/gen/go/category_service/v1"
	"github.com/Kazzess/libraries/apperror"
	"github.com/Kazzess/libraries/core/closer"
	"github.com/Kazzess/libraries/core/healthcheck"
	"github.com/Kazzess/libraries/core/safe"
	"github.com/Kazzess/libraries/errors"
	"github.com/Kazzess/libraries/logging"
	"github.com/Kazzess/libraries/metrics"
	psql "github.com/Kazzess/libraries/postgresql"
	"github.com/Kazzess/libraries/tracing"
	"github.com/Kazzess/libraries/utils/clock"
	"github.com/Kazzess/libraries/utils/ident"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"crm/app/internal/config"
	gRPCCategory "crm/app/internal/controller/grpc/v1/category"
	"crm/app/internal/dal/postgres"
	"crm/app/internal/domain"
	domainCategoryService "crm/app/internal/domain/category/service"
	domainCategoryStorage "crm/app/internal/domain/category/storage/postgres"
	"crm/app/internal/policy"
	policyCategory "crm/app/internal/policy/category"
)

type Runner interface {
	Run(context.Context) error
}

type App struct {
	cfg        *config.Config
	gRPCServer *grpc.Server
	httpRouter *chi.Mux
	httpServer *http.Server

	metricsHTTTPServer *metrics.Server
	healthServer       *healthcheck.GRPCHealthServer

	policyCategory *policyCategory.Policy

	runners []Runner
}

func (a *App) AddRunner(runner Runner) {
	a.runners = append(a.runners, runner)
}

//nolint:funlen
func NewApp(ctx context.Context) (*App, error) {
	app := App{}

	cfg := config.GetConfig()
	app.cfg = cfg

	logger := logging.NewLogger(
		logging.WithLevel(cfg.App.LogLevel),
		logging.WithIsJSON(cfg.App.IsLogJSON),
	)
	ctx = logging.ContextWithLogger(ctx, logger)

	logging.L(ctx).Info("config loaded", "config", cfg)

	app.healthServer = healthcheck.NewGRPCHealthServer()

	// Init Postgres Client.
	postgresClient, err := app.initPostgresClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't create postgres Client")
	}

	uuidGenerator := ident.NewUUIDGenerator()
	defClock := clock.NewDefault()

	// Init storage and service.
	categoryStorage := domainCategoryStorage.NewStorage(postgresClient)
	categoryService := domainCategoryService.NewService(categoryStorage)

	// Init policy.
	basePolicy := policy.NewBasePolicy(
		uuidGenerator,
		defClock,
	)

	app.policyCategory = policyCategory.NewPolicy(
		basePolicy,
		categoryService,
	)

	// init gRPC controllers
	app.gRPCServer = app.initGRPCServer(ctx)

	// init HTTP router
	app.httpRouter = app.initHTTPRouter(ctx)

	return &app, nil
}

func (a *App) Run(ctx context.Context) error {
	// Run migrations.
	err := postgres.RunMigrations(&a.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "migrations failed")
	}

	errGroup, _ := safe.WithContext(ctx)

	errGroup.Run(closer.CloseOnSignalContext(os.Kill, os.Interrupt))
	errGroup.Run(a.setupGRPCServer)
	errGroup.Run(a.setupHTTPServer)

	for _, r := range a.runners {
		errGroup.Run(r.Run)
	}

	logging.L(ctx).Info("application started")

	return errGroup.Wait()
}

func (a *App) setupGRPCServer(ctx context.Context) error {
	logging.L(ctx).Info(
		"gRPC server initializing",
		logging.StringAttr("host", a.cfg.GRPC.Host),
		logging.IntAttr("port", a.cfg.GRPC.Port),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port))
	if err != nil {
		return errors.Wrap(err, "gRPC server listen error")
	}

	closer.Add(lis)

	if err = a.gRPCServer.Serve(lis); err != nil {
		return errors.Wrap(err, "gRPC server serve error")
	}

	return nil
}

func (a *App) setupHTTPServer(ctx context.Context) error {
	logging.L(ctx).Info(
		"HTTP server initializing",
		logging.StringAttr("host", a.cfg.HTTP.Host),
		logging.IntAttr("port", a.cfg.HTTP.Port),
		logging.DurationAttr("read_timeout", a.cfg.HTTP.ReadHeaderTimeout),
	)

	a.httpServer = &http.Server{
		Addr:        fmt.Sprintf("%s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port),
		Handler:     a.httpRouter,
		ReadTimeout: a.cfg.HTTP.ReadHeaderTimeout,
	}

	closer.Add(a.httpServer)

	if err := a.httpServer.ListenAndServe(); err != nil {
		logging.L(ctx).With(logging.ErrAttr(err)).Error("HTTP server listen and serve error")
		return err
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) *grpc.Server {
	logging.L(ctx).Info(
		"gRPC server initializing",
		logging.StringAttr("host", a.cfg.GRPC.Host),
		logging.IntAttr("port", a.cfg.GRPC.Port),
	)

	recoveryHandler := grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "internal system error")
	})

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			apperror.GRPCUnaryInterceptor(domain.SystemCode),
			logging.WithTraceIDInLogger(),
			metrics.RequestDurationMetricUnaryServerInterceptor(fmt.Sprintf(
				"%s-%s-%s",
				a.cfg.App.Name,
				a.cfg.App.ID,
				a.cfg.App.Version,
			)),
			grpc_recovery.UnaryServerInterceptor(recoveryHandler),
		),
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(recoveryHandler),
		),
	}

	serverOptions = append(serverOptions, tracing.WithAllTracing()...)

	gRPCServer := grpc.NewServer(
		serverOptions...,
	)
	reflection.Register(gRPCServer)

	grpc_health_v1.RegisterHealthServer(gRPCServer, a.healthServer.Server)

	a.healthServer.HealthCheck(ctx, a.cfg.GRPC.HealthCheckInterval)

	// Init controllers.
	gRPCCategoryService.RegisterCategoryServiceServer(gRPCServer,
		gRPCCategory.NewController(
			a.policyCategory,
		),
	)

	return gRPCServer
}

func (a *App) initHTTPRouter(_ context.Context) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logging.Middleware)
	router.Use(tracing.Middleware)

	router.Use(metrics.RequestDurationMetricHTTPMiddleware)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(config.Timeout))

	return router
}

func (a *App) initPostgresClient(ctx context.Context) (*psql.Client, error) {
	logging.WithAttrs(
		ctx,
		logging.StringAttr("host", a.cfg.Postgres.Host),
		logging.IntAttr("port", a.cfg.Postgres.Port),
		logging.StringAttr("user", a.cfg.Postgres.User),
		logging.StringAttr("db", a.cfg.Postgres.Database),
		logging.StringAttr("password", "<REMOVED>"),
		logging.IntAttr("max-attempts", a.cfg.Postgres.MaxAttempt),
		logging.DurationAttr("max_delay", a.cfg.Postgres.MaxDelay),
	).Info("PostgreSQL initializing")

	pgDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		a.cfg.Postgres.User,
		a.cfg.Postgres.Password,
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.Database,
	)

	if a.cfg.Postgres.Binary {
		pgDsn += "?sslmode=require"
	}

	postgresConfig, err := psql.NewConfig(
		pgDsn,
		a.cfg.Postgres.MaxAttempt,
		a.cfg.Postgres.MaxDelay,
		psql.WithBinaryExecMode(a.cfg.Postgres.Binary),
		psql.WithHealthChecker("", a.healthServer),
	)
	if err != nil {
		return nil, errors.Wrap(err, "psql.NewConfig")
	}

	pgClient, err := psql.NewClient(ctx, postgresConfig)
	if err != nil {
		return nil, errors.Wrap(err, "psql.NewClient")
	}

	closer.AddN(pgClient)

	return pgClient, nil
}
