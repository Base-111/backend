package http

import (
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DeleteHandler struct {
	uc prompt.DeletePromptUC
}

func NewDeleteHandler(uc prompt.DeletePromptUC) *DeleteHandler {
	return &DeleteHandler{
		uc: uc,
	}
}

func (h *DeleteHandler) Handle(c *gin.Context) {
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

	err = h.uc.Execute(c, parsedId)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})

		return
	}

	c.Status(http.StatusNoContent)
}
