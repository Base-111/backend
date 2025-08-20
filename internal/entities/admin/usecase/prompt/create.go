package prompt

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/repository"
	"github.com/gin-gonic/gin"
)

type CreatePromptUC interface {
	Execute(ctx *gin.Context, product domain.Prompt) error
}

type CreatePromptUseCase struct {
	repository repository.PromptRepo
}

func NewCreatePromptUseCase(repo repository.PromptRepo) *CreatePromptUseCase {
	return &CreatePromptUseCase{repository: repo}
}

func (uc *CreatePromptUseCase) Execute(ctx *gin.Context, product domain.Prompt) error {
	return uc.repository.Insert(ctx, product)
}
