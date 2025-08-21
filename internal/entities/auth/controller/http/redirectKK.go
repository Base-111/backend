package http

import (
	"github.com/Base-111/backend/internal/utils"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RedirectToKeycloak godoc
// @Summary      Начало авторизации
// @Description  Перенаправляет пользователя на провайдера (Keycloak) для авторизации. Возвращает 302 с заголовком Location.
// @Tags         auth
// @Produce      json
// @Success      302  {string}  string  "Found"
// @Failure      500  {string}  string  "internal error"
// @Router       /auth/login [get]
func (a *AuthHandler) RedirectToKeycloak(c *gin.Context) {
	stateID, err := utils.GenerateRandomBase64Str()
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
	}
	if err = a.authStore.SetState(c, stateID); err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})
	}
	c.Redirect(http.StatusFound, a.authClient.Oauth.AuthCodeURL(stateID))
}
