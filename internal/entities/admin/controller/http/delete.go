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

// Handle godoc
// @Summary      Удалить prompt
// @Description  Удаляет prompt по ID.
// @Tags         prompt
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      204  {string}  string  "no content"
// @Failure      400  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Router       /admin/prompt/{id} [delete]
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
