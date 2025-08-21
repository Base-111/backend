package router

import (
	_ "github.com/Base-111/backend/docs"
	admin "github.com/Base-111/backend/internal/entities/admin/controller/http/router"
	auth "github.com/Base-111/backend/internal/entities/auth/controller/http"
	authRouter "github.com/Base-111/backend/internal/entities/auth/controller/http/router"
	"github.com/Base-111/backend/pkg/logs/api"
	"github.com/Base-111/backend/pkg/tracing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		cors.Default(),
		api.RequestLoggingMiddleware(),
		api.LoggerMiddleware(),
		tracing.Middleware(),
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	authRouter.InitAuthRoutes(router, h.authHandler)
	admin.InitAdminRoutes(router, h.adminHandler)

	return router, nil
}
