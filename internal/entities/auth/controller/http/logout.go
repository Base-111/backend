package http

import (
	"fmt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *AuthHandler) Logout(c *gin.Context) {
	sessionId, err := c.Cookie("session_id")
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
	}

	c.SetCookie(
		"session_id", // name
		"",           // value (user ID)
		3600,         // maxAge (1 hour)
		"/",          // path
		"",           // domain
		true,         // secure
		true,         // httpOnly
	)

	c.SetCookie(
		"user_email",
		"",
		3600,
		"/",
		"",
		true,
		true,
	)
	err = a.sessionStore.DeleteSession(c.Request.Context(), sessionId)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", a.cfg.Auth.BaseURL, a.cfg.Auth.Realm))
}
