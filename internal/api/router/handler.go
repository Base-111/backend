package router

import (
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (h *ApiHandler) SetupRoutes() (*gin.Engine, error) {
	router := gin.New()

	router.Use(
		gin.Recovery(),
	)

	return router, nil
}
