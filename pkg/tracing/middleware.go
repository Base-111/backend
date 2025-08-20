package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const Header = "X-Request-Id"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := uuid.Parse(c.Request.Header.Get(Header))
		if requestID == uuid.Nil {
			requestID = uuid.New()
		}

		ctx := WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
