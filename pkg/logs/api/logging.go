package api

import (
	"fmt"
	"github.com/Base-111/backend/pkg/logs"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net"
	"time"
)

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			host = c.Request.RemoteAddr
		}

		logger := logs.Logger(c.Request.Context()).With(
			slog.String("request_user_agent", c.Request.UserAgent()),
			slog.String("request_referrer", c.Request.Referer()),
			slog.String("request_host", host),
		)

		message := fmt.Sprintf("request %s %s", c.Request.Method, c.Request.URL.Path)

		logger.Info(message + " started")

		defer func() {
			logger.
				With(slog.Int64("request_latency_ms", time.Since(startAt).Milliseconds())).
				Info(message + " completed")
		}()

		c.Next()
	}
}
