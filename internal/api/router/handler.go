package router

import (
	admin "github.com/Base-111/backend/internal/entities/admin/controller/http/router"
	auth "github.com/Base-111/backend/internal/entities/auth/controller/http"
	authRouter "github.com/Base-111/backend/internal/entities/auth/controller/http/router"
	"github.com/Base-111/backend/pkg/logs/api"
	"github.com/Base-111/backend/pkg/tracing"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	authHandler  *auth.AuthHandler
	adminHandler *admin.HandlerContainer
}

func NewApiHandler(authHandler *auth.AuthHandler, adminHandler *admin.HandlerContainer) *ApiHandler {
	return &ApiHandler{
		authHandler:  authHandler,
		adminHandler: adminHandler,
	}
}

func (h *ApiHandler) SetupRoutes() (*gin.Engine, error) {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		api.RequestLoggingMiddleware(),
		api.LoggerMiddleware(),
		tracing.Middleware(),
	)

	authRouter.InitAuthRoutes(router, h.authHandler)
	admin.InitAdminRoutes(router, h.adminHandler)

	return router, nil
}
