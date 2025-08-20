package http

import (
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReadHandler struct {
	uc prompt.ReadPromptUC
}

func NewReadHandler(uc prompt.ReadPromptUC) *ReadHandler {
	return &ReadHandler{
		uc: uc,
	}
}

func (h *ReadHandler) Handle(c *gin.Context) {
	id := c.Param("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewParseError(err),
		})

		return
	}

	order, err := h.uc.Execute(c, parsedId)

	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})

		return
	}

	c.JSON(http.StatusOK, order)
}
