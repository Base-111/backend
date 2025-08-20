package http

import (
	"github.com/Base-111/backend/internal/config"
	"github.com/Base-111/backend/internal/entities/auth/store"
	"github.com/Base-111/backend/pkg/auth"
)

type AuthHandler struct {
	cfg          *config.Config
	serverAddr   string
	authClient   *auth.Client
	authStore    store.AuthStore
	sessionStore store.SessionStore
}

func New(
	cfg *config.Config,
	serverAddr string,
	authClient *auth.Client,
	authStore store.AuthStore,
	sessionStore store.SessionStore,
) *AuthHandler {
	return &AuthHandler{
		cfg:          cfg,
		serverAddr:   serverAddr,
		authClient:   authClient,
		authStore:    authStore,
		sessionStore: sessionStore,
	}
}
