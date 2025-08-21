package middleware

import (
	"github.com/Base-111/backend/internal/entities/auth/store"
	"github.com/Base-111/backend/pkg/auth"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type TokenClaims struct {
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
	Name     string `json:"name"`

	Scope       string `json:"scope"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

func AuthMiddleware(sessionStore store.SessionStore, authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userEmail, err := c.Cookie("user_email")
		if err != nil || userEmail == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		sessionData, err := sessionStore.GetSession(c, sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session",
			})
			c.Abort()
			return
		}
		if sessionData.AccessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session not found",
			})
			c.Abort()
			return
		}
		if authClient.Provider == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "OIDC provider not initialized",
			})
			c.Abort()
			return
		}

		keySet := authClient.Provider.VerifierContext(c, &oidc.Config{
			SkipClientIDCheck: true,
		})
		if keySet == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid to verify provider access token",
			})
			c.Abort()
			return
		}

		token, err := keySet.Verify(c, sessionData.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid verify access token",
			})
			c.Abort()
			return
		}

		var claims TokenClaims
		if err := token.Claims(&claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid parse claims",
			})
			c.Abort()
			return
		}

		c.Set("user_id", sessionID)
		c.Set("user_email", userEmail)

		c.Next()
	}
}
