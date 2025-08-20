package sql

import (
	"context"
	"fmt"
	"github.com/Base-111/backend/pkg/logs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/pkg/errors"
	"log/slog"
	"regexp"
)

const (
	maxArgs        = 30
	maxQueryLength = 1_000
)

type ContextLogger struct{}

func (log *ContextLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	if err, exists := data["err"].(error); exists {
		data["err"] = fmt.Sprintf("%+v", err)
		if isWarningError(err) {
			level = tracelog.LogLevelWarn
		}
	}

	fields, msg := log.formatFieldsAndMessage(data, msg)
	logger := logs.Logger(ctx).With(fields...)

	logMessage(logger, msg, level)
}

func (log *ContextLogger) formatFieldsAndMessage(data map[string]any, msg string) ([]any, string) {
	fields := make([]any, 0, len(data)+1)

	for key, value := range data {
		switch key {
		case "args":
			args, ok := value.([]any)
			if ok && len(args) > 0 && len(args) <= maxArgs {
				fields = append(fields, slog.Any(key, value))
			}
		case "sql":
			sql, ok := value.(string)
			if ok {
				shortSQL := formatLoggedSQL(sql)
				if shortSQL != "" {
					msg += ": " + shortSQL
				}
				fields = append(fields, slog.String(key, cutString(sql)))
			} else {
				fields = append(fields, slog.Any(key, value))
			}
		default:
			fields = append(fields, slog.Any(key, value))
		}
	}

	return fields, msg
}

func logMessage(logger *slog.Logger, msg string, level tracelog.LogLevel) {
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		logger.Debug(msg)
	case tracelog.LogLevelInfo:
		logger.Info(msg)
	case tracelog.LogLevelWarn:
		logger.Warn(msg)
	case tracelog.LogLevelError:
		logger.Error(msg)
	default:
		logger.With("INVALID_PGX_LOG_LEVEL", level).Error(msg)
	}
}

func isWarningError(err error) bool {
	return pgconn.Timeout(err) || errors.Is(err, context.Canceled)
}

type sqlParser struct {
	Matcher *regexp.Regexp
	Format  string
}

var sqlParsers = []sqlParser{
	{
		Matcher: regexp.MustCompile(`(?i)SELECT\s+.*\s+FROM (?P<tableName>[a-zA-Z_]+(\.[a-zA-Z_]+)?).*`),
		Format:  "SELECT ... FROM %s",
	},
	{
		Matcher: regexp.MustCompile(`(?i)INSERT\s*.*\s*INTO (?P<tableName>[a-zA-Z_]+(\.[a-zA-Z_]+)?).*`),
		Format:  "INSERT INTO %s ...",
	},
	{
		Matcher: regexp.MustCompile(`(?i)UPDATE\s+(?P<tableName>[a-zA-Z_]+(\.[a-zA-Z_]+)?).*`),
		Format:  "UPDATE %s ...",
	},
	{
		Matcher: regexp.MustCompile(`(?i)DELETE\s*.*\s*FROM (?P<tableName>[a-zA-Z_]+(\.[a-zA-Z_]+)?).*`),
		Format:  "DELETE FROM %s ...",
	},
	{
		Matcher: regexp.MustCompile(`(?i)(?P<operation>[a-zA-Z]+)\s*.*`),
		Format:  "%s",
	},
}

func formatLoggedSQL(sql string) string {
	for _, parser := range sqlParsers {
		matches := parser.Matcher.FindStringSubmatch(sql)
		if len(matches) > 1 {
			return fmt.Sprintf(parser.Format, matches[1])
		}
	}

	return ""
}

func cutString(s string) string {
	if len(s) <= maxQueryLength {
		return s
	}

	return s[:maxQueryLength] + "..."
}
