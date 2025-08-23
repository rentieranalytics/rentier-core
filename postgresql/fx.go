package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"postgresql",
	fx.Provide(
		func() Options {
			return Options{
				CaptureDBErrors: true,
			}
		},
		InitTracer,
		InitPostgresqlPool,
	),
)

type PostgresqlConfigurer interface {
	GetPostgresqlServerAddress() string
}

func InitTracer(opt Options) pgx.QueryTracer {
	return NewTracer(opt)
}

func InitPostgresqlPool(
	config PostgresqlConfigurer,
	tracer pgx.QueryTracer,
) (*pgxpool.Pool, error) {
	ctx := context.Background()

	dbConfig, err := pgxpool.ParseConfig(config.GetPostgresqlServerAddress())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	dbConfig.MaxConns = 20
	dbConfig.MinConns = 2
	dbConfig.MaxConnIdleTime = 5 * time.Minute
	dbConfig.MaxConnLifetime = 30 * time.Minute
	dbConfig.MaxConnLifetimeJitter = 5 * time.Minute

	dbConfig.HealthCheckPeriod = 1 * time.Minute
	dbConfig.ConnConfig.ConnectTimeout = 5 * time.Second
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	dbConfig.ConnConfig.Tracer = multitracer.New(tracer)
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}
