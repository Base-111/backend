package repository

import (
	"context"
	"github.com/Base-111/backend/internal/entities/admin/domain"
)

type PromptRepo interface {
	Insert(ctx context.Context, Prompt domain.Prompt) error
	GetById(ctx context.Context, PromptId int) (domain.Prompt, error)
	GetAll(ctx context.Context, params domain.PromptFilterParams) ([]domain.Prompt, error)
	Update(ctx context.Context, PromptId int, Prompt domain.Prompt) error
	Delete(ctx context.Context, PromptId int) error
}
