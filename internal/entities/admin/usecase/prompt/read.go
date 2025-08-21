package prompt

import (
	"github.com/Base-111/backend/internal/entities/admin/domain"
	"github.com/Base-111/backend/internal/entities/admin/repository"
	"github.com/gin-gonic/gin"
)

type ReadPromptUC interface {
	Execute(ctx *gin.Context, id int) (domain.Prompt, error)
	ExecuteAll(ctx *gin.Context, params domain.PromptFilterParams) ([]domain.Prompt, error)
}

type ReadPromptUseCase struct {
	repository repository.PromptRepo
}

func NewReadPromptUseCase(repo repository.PromptRepo) *ReadPromptUseCase {
	return &ReadPromptUseCase{repository: repo}
}

func (uc *ReadPromptUseCase) Execute(ctx *gin.Context, id int) (domain.Prompt, error) {
	return uc.repository.GetById(ctx, id)
}

func (uc *ReadPromptUseCase) ExecuteAll(ctx *gin.Context, params domain.PromptFilterParams) ([]domain.Prompt, error) {
	return uc.repository.GetAll(ctx, params)
}
