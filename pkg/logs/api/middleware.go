package api

import (
	"github.com/Base-111/backend/pkg/logs"
	"github.com/Base-111/backend/pkg/tracing"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := tracing.RequestID(c.Request.Context())

		logger := slog.Default().With(
			slog.String("request_id", requestID.String()),
			slog.String("request_method", c.Request.Method),
			slog.String("request_url", c.Request.URL.Path),
		)

		ctx := logs.WithLogger(c.Request.Context(), logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
