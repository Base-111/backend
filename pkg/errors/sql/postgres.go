package sql

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"log/slog"
	"slices"
)

type Error struct {
	SQL  string
	Args []any
	Err  error

	attrs []slog.Attr
}

func NewSQLError(sql string, args []any, err error) error {
	return errors.WithStack(&Error{
		SQL:   sql,
		Args:  args,
		Err:   err,
		attrs: getPgErrorAttrs(err),
	})
}

func WrapSQLError(message, sql string, args []any, err error) error {
	return errors.Wrap(
		&Error{
			SQL:   sql,
			Args:  args,
			Err:   err,
			attrs: getPgErrorAttrs(err),
		},
		message,
	)
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) LogAttrs() []slog.Attr {
	return slices.Concat(
		[]slog.Attr{
			slog.String("sql", e.SQL),
			slog.Any("args", e.Args),
		},
		e.attrs,
	)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func getPgErrorAttrs(err error) []slog.Attr {
	var e *pgconn.PgError
	if !errors.As(err, &e) {
		return nil
	}

	options := make([]slog.Attr, 0, 18)

	if e.Severity != "" {
		options = append(options, slog.String("pg_severity", e.Severity))
	}
	if e.Code != "" {
		options = append(options, slog.String("pg_code", e.Code))
	}
	if e.Message != "" {
		options = append(options, slog.String("pg_message", e.Message))
	}
	if e.Detail != "" {
		options = append(options, slog.String("pg_detail", e.Detail))
	}
	if e.Hint != "" {
		options = append(options, slog.String("pg_hint", e.Hint))
	}
	if e.Where != "" {
		options = append(options, slog.String("pg_where", e.Where))
	}
	if e.SchemaName != "" {
		options = append(options, slog.String("pg_schema_name", e.SchemaName))
	}
	if e.TableName != "" {
		options = append(options, slog.String("pg_table_name", e.TableName))
	}
	if e.ColumnName != "" {
		options = append(options, slog.String("pg_column_name", e.ColumnName))
	}
	if e.DataTypeName != "" {
		options = append(options, slog.String("pg_data_type_name", e.DataTypeName))
	}
	if e.ConstraintName != "" {
		options = append(options, slog.String("pg_constraint_name", e.ConstraintName))
	}

	return options
}
