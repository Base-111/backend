package prompt

import (
	"github.com/Base-111/backend/internal/entities/admin/repository"
	"github.com/gin-gonic/gin"
)

type DeletePromptUC interface {
	Execute(ctx *gin.Context, id int) error
}

type DeletePromptUseCase struct {
	repository repository.PromptRepo
}

func NewDeletePromptUseCase(repo repository.PromptRepo) *DeletePromptUseCase {
	return &DeletePromptUseCase{repository: repo}
}

func (uc *DeletePromptUseCase) Execute(ctx *gin.Context, id int) error {
	return uc.repository.Delete(ctx, id)
}
