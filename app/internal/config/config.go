package config

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	ID            string `yaml:"id" env:"APP_ID"`
	Name          string `yaml:"name" env:"APP_NAME"`
	Version       string `yaml:"version" env:"APP_VERSION"`
	IsDevelopment bool   `yaml:"is_dev" env:"APP_IS_DEVELOPMENT"`
	LogLevel      string `yaml:"log_level" env:"APP_LOG_LEVEL"`
	IsLogJSON     bool   `yaml:"is_log_json" env:"APP_IS_LOG_JSON"`
	Domain        string `yaml:"domain" env:"APP_DOMAIN"`
}

type GRPCConfig struct {
	Host                string        `yaml:"host" env:"GRPC_HOST"`
	Port                int           `yaml:"port" env:"GRPC_PORT"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval" env:"GRPC_HEALTH_CHECK_INTERVAL"`
}

type HTTPConfig struct {
	Host              string        `yaml:"host" env:"HTTP_HOST"`
	Port              int           `yaml:"port" env:"HTTP_PORT"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
}

type PostgresConfig struct {
	Host       string        `yaml:"host"  env:"POSTGRES_HOST"`
	User       string        `yaml:"user" env:"POSTGRES_USER"`
	Password   string        `yaml:"password" env:"POSTGRES_PASSWORD"`
	Port       int           `yaml:"port" env:"POSTGRES_PORT"`
	Database   string        `yaml:"database" env:"POSTGRES_DATABASE"`
	MaxAttempt int           `yaml:"max_attempt"`
	MaxDelay   time.Duration `yaml:"max_delay"`
	Binary     bool          `yaml:"binary" env:"POSTGRES_BINARY"`
}

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Postgres PostgresConfig `yaml:"postgres"`
	HTTP     HTTPConfig     `yaml:"http"`
}

const (
	FlagConfigPathName = "config"
	EnvConfigPathName  = "CONFIG_PATH"
)

var (
	configPath string
	instance   Config
	once       sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(
			&configPath,
			FlagConfigPathName,
			"/home/kazess/workspace/crm/configs/config.local.yaml",
			"this is application configuration file",
		)
		flag.Parse()

		if path, ok := os.LookupEnv(EnvConfigPathName); ok {
			configPath = path
		}

		log.Printf("config initializing from: %s", configPath)

		instance = Config{}

		if err := cleanenv.ReadConfig(configPath, &instance); err != nil {
			help, _ := cleanenv.GetDescription(&instance, nil)
			log.Println(help)
			log.Fatal(err)
		}

		log.Println("configuration loaded")
	})

	return &instance
}
