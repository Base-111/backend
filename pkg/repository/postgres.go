package repository

import (
	"context"
	"fmt"
	"github.com/Base-111/backend/pkg/logs/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"
)

type Config struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func GetQueryBuilderFormat() sq.StatementBuilderType {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return qb
}

func ConnectViaPGXConnect(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	databaseUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, errors.Wrap(err, "parse postgres config")
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &sql.ContextLogger{},
		LogLevel: tracelog.LogLevelInfo,
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, errors.Wrap(err, "create postgres pool")
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "ping postgres")
	}

	return pool, nil
}
