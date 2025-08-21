package http

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReadAllHandler struct {
	uc prompt.ReadPromptUC
}

func NewReadAllHandler(uc prompt.ReadPromptUC) *ReadAllHandler {
	return &ReadAllHandler{
		uc: uc,
	}
}

// Handle godoc
// @Summary      Список prompts
// @Description  Возвращает список prompts с пагинацией.
// @Tags         prompt
// @Produce      json
// @Param        page       query     int   false  "Номер страницы"       minimum(1)
// @Param        page_size  query     int   false  "Размер страницы"       minimum(1) maximum(100)
// @Success      200        {array}   domain.Prompt
// @Failure      400        {string}  string
// @Failure      500        {string}  string
// @Router       /admin/prompt/ [get]
func (h *ReadAllHandler) Handle(c *gin.Context) {
	pageParam := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pageParam, 10, 64)

	if err != nil {
		page = 1
	}

	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitParam, 10, 64)

	if err != nil {
		limit = 10
	}

	params := domain.PromptFilterParams{
		Page:     page,
		PageSize: limit,
	}

	products, err := h.uc.ExecuteAll(c, params)

	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})

		return
	}

	c.JSON(http.StatusOK, products)
}
