package tracing

import (
	"context"
	"github.com/google/uuid"
)

type traceKey struct{}

func RequestID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(traceKey{}).(uuid.UUID); ok {
		return id
	}

	return uuid.New()
}

func WithRequestID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, traceKey{}, id)
}
