package router

import (
	"github.com/Base-111/backend/internal/entities/admin/controller/http"
	"github.com/gin-gonic/gin"
)

type HandlerContainer struct {
	CreateHandler  *http.CreateHandler
	ReadHandler    *http.ReadHandler
	ReadAllHandler *http.ReadAllHandler
	DeleteHandler  *http.DeleteHandler
	UpdateHandler  *http.UpdateHandler
}

func InitAdminRoutes(router *gin.Engine, handler *HandlerContainer) {
	routerGroup := router.Group("/admin")
	{
		promptRouter := routerGroup.Group("/prompt")
		{
			promptRouter.GET("/", handler.ReadAllHandler.Handle)
			promptRouter.GET("/:id", handler.ReadHandler.Handle)
			promptRouter.POST("/", handler.CreateHandler.Handle)
			promptRouter.PUT("/:id", handler.UpdateHandler.Handle)
			promptRouter.DELETE("/:id", handler.DeleteHandler.Handle)
		}
	}
}
