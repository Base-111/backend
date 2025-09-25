package router

import (
	_ "github.com/Base-111/backend/docs"
	auth "github.com/Base-111/backend/internal/api/middleware"
	admin "github.com/Base-111/backend/internal/entities/admin/controller/http/router"
	"github.com/Base-111/backend/pkg/logs/api"
	"github.com/Base-111/backend/pkg/tracing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ApiHandler struct {
	adminHandler *admin.HandlerContainer
}

func NewApiHandler(adminHandler *admin.HandlerContainer) *ApiHandler {
	return &ApiHandler{
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
		auth.AuthMiddleware(),
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	admin.InitAdminRoutes(router, h.adminHandler)

	return router, nil
}
