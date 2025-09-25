package middleware

import (
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

var Verifier *oidc.IDTokenVerifier

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}
		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := Verifier.Verify(c.Request.Context(), rawToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			c.Abort()
			return
		}

		var claims map[string]interface{}
		if err = idToken.Claims(&claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
