package logs

import (
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

type Loggable interface {
	LogAttrs() []slog.Attr
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func Error(request *http.Request, err error) {
	logger := Logger(request.Context())

	for e := err; e != nil; e = errors.Unwrap(e) {
		if loggable, ok := e.(Loggable); ok {
			logger = loggerWithAttrs(logger, loggable.LogAttrs())
		}
		if trace, ok := e.(StackTracer); ok {
			logger = logger.With("stacktrace", trace.StackTrace())
		}
	}

	logger.Error(fmt.Sprintf("handle %s %s: %v", request.Method, request.URL.Path, err))
}

func loggerWithAttrs(logger *slog.Logger, attrs []slog.Attr) *slog.Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}
