package router

import (
	"github.com/Base-111/backend/internal/entities/auth/controller/http"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(router *gin.Engine, handler *http.AuthHandler) {
	routerGroup := router.Group("/auth")
	{
		routerGroup.GET("/login", handler.RedirectToKeycloak)
		routerGroup.GET("/callback", handler.Callback)
		routerGroup.GET("/logout", handler.Logout)
	}
}
