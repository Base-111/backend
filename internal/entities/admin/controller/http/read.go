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

// Handle godoc
// @Summary      Получить prompt
// @Description  Возвращает prompt по ID.
// @Tags         prompt
// @Produce      json
// @Param        id   path      int  true  "ID"
// @Success      200  {object}  domain.Prompt
// @Failure      400  {string}  string
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Router       /admin/prompt/{id} [get]
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
