package prompt

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/repository"
	"github.com/gin-gonic/gin"
)

type UpdatePromptUC interface {
	Execute(ctx *gin.Context, id int, product domain.Prompt) error
}

type UpdatePromptUseCase struct {
	repository repository.PromptRepo
}

func NewUpdatePromptUseCase(repo repository.PromptRepo) *UpdatePromptUseCase {
	return &UpdatePromptUseCase{repository: repo}
}

func (uc *UpdatePromptUseCase) Execute(ctx *gin.Context, id int, prompt domain.Prompt) error {
	return uc.repository.Update(ctx, id, prompt)
}
