package http

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/errors/api"
	"github.com/gin-gonic/gin"
	"github.com/muonsoft/validation/validator"
	"net/http"
)

type CreateHandler struct {
	uc prompt.CreatePromptUC
}

func NewCreateHandler(uc prompt.CreatePromptUC) *CreateHandler {
	return &CreateHandler{
		uc: uc,
	}
}

// Handle godoc
// @Summary      Создать prompt
// @Description  Создаёт новый prompt.
// @Tags         prompt
// @Accept       json
// @Produce      json
// @Param        prompt  body      domain.Prompt  true  "Новый prompt"
// @Success      201     {string}  string         "created"
// @Failure      400     {string}  string
// @Failure      500     {string}  string
// @Router       /admin/prompt/ [post]
func (h *CreateHandler) Handle(c *gin.Context) {
	var inputJson domain.Prompt
	if err := c.ShouldBindJSON(&inputJson); err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewUnmarshalError(err),
		})

		return
	}

	err := validator.ValidateIt(c.Request.Context(), &inputJson)
	if err != nil {
		api.WriteError(c, &api.HandlerError{
			Method: c.Request.Method,
			URL:    c.Request.URL.Path,
			Err:    api.NewValidationError(err),
		})

		return
	}

	err = h.uc.Execute(c, inputJson)
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
