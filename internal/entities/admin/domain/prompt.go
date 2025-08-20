package domain

import (
	"context"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type Prompt struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Text   string `json:"text,omitempty"`
	System bool   `json:"is_system,omitempty"`
}

type PromptFilterParams struct {
	Page     int64
	PageSize int64
}

func (p *Prompt) Validate(ctx context.Context, v *validation.Validator) error {
	return v.Validate(
		ctx,
		validation.StringProperty("name", p.Name, it.IsNotBlank(), it.HasMinLength(3)),
	)
}
