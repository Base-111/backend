package http

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"github.com/muonsoft/validation/validator"
	"net/http"
	"strconv"
)

type UpdateHandler struct {
	uc prompt.UpdatePromptUC
}

func NewUpdateHandler(uc prompt.UpdatePromptUC) *UpdateHandler {
	return &UpdateHandler{
		uc: uc,
	}
}

// Handle godoc
// @Summary      Обновить prompt
// @Description  Обновляет существующий prompt по ID.
// @Tags         prompt
// @Accept       json
// @Produce      json
// @Param        id      path      int           true  "ID"
// @Param        prompt  body      domain.Prompt true  "Данные для обновления"
// @Success      201     {string}  string        "updated"
// @Failure      400     {string}  string
// @Failure      404     {string}  string
// @Failure      500     {string}  string
// @Router       /admin/prompt/{id} [put]
func (h *UpdateHandler) Handle(c *gin.Context) {
	var input domain.Prompt
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewParseError(err),
		})

		return
	}
	if err = c.BindJSON(&input); err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewUnmarshalError(err),
		})

		return
	}

	err = validator.ValidateIt(c.Request.Context(), &input)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewValidationError(err),
		})

		return
	}

	err = h.uc.Execute(c, id, input)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    err,
		})

		return
	}

	c.JSON(http.StatusCreated, nil)
}
