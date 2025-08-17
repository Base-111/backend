package http

import (
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

func (h *ApiHandler) InitHandlers() *gin.Engine {
	gin.SetMode(gin.TestMode)

	return nil
}
