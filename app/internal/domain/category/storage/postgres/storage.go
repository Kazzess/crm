package postgres

import (
	"context"
	"strconv"

	psql "github.com/Kazzess/libraries/postgresql"
	"github.com/Kazzess/libraries/tracing"
	"github.com/Masterminds/squirrel"

	"crm/app/internal/dal/postgres"
	domainCategory "crm/app/internal/domain/category"
	modelCategory "crm/app/internal/domain/category/model"
)

type Storage struct {
	qb     squirrel.StatementBuilderType
	client *psql.Client
}

func NewStorage(client *psql.Client) *Storage {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &Storage{client: client, qb: qb}
}

func (repo *Storage) Create(ctx context.Context, res modelCategory.CreateCategoryInput) error {
	query, args, err := repo.qb.
		Insert(postgres.CategoryTable.String()).
		Columns(
			"name",
			"is_enabled",
			"created_at",
			"updated_at",
		).
		Values(
			res.Name,
			res.IsEnabled,
			res.CreatedAt,
			res.UpdatedAt,
		).
		ToSql()
	if err != nil {
		err = psql.ErrCreateQuery(err)
		tracing.Error(ctx, err)

		return err
	}

	tracing.SpanEvent(ctx, "create Category query")
	tracing.TraceValue(ctx, "sql", query)

	for i, arg := range args {
		tracing.TraceValue(ctx, strconv.Itoa(i), arg)
	}

	_, execErr := repo.client.Exec(ctx, query, args...)
	if execErr != nil {
		if pgErr, ok := psql.IsErrUniqueViolation(execErr); ok {
			switch pgErr.ConstraintName {
			case domainCategory.CategoryIDPkConstraint:
				return domainCategory.ErrViolatesConstraintCategoryIDPk
			}

			execErr = psql.ErrDoQuery(psql.ParsePgError(execErr))

			return execErr
		}
	}

	return nil
}
