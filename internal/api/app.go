package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Base-111/backend/internal/api/middleware"
	"github.com/Base-111/backend/internal/api/router"
	"github.com/Base-111/backend/internal/api/server"
	"github.com/Base-111/backend/internal/config"
	"github.com/Base-111/backend/internal/entities/admin/controller/http"
	admin "github.com/Base-111/backend/internal/entities/admin/controller/http/router"
	"github.com/Base-111/backend/internal/entities/admin/repository/postgres"
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	"github.com/Base-111/backend/pkg/repository"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message/translations/english"
	"github.com/muonsoft/validation/validator"
	"golang.org/x/text/language"
)

func Run() error {
	ctx := context.Background()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}
	err = validator.SetUp(
		validation.DefaultLanguage(language.English),
		validation.Translations(english.Messages),
	)

	if err != nil {
		return err
	}

	serv := new(server.Server)

	keycloakURL := "http://localhost:8080/realms/myrealm"
	provider, err := oidc.NewProvider(ctx, keycloakURL)
	if err != nil {
		return err
	}

	middleware.Verifier = provider.Verifier(&oidc.Config{
		ClientID: "myclient",
	})

	dbConnection, err := repository.ConnectViaPGXConnect(ctx, cfg.Postgres)

	if err != nil {
		return err
	}

	promptRepo := postgres.NewPromptRepository(dbConnection, repository.GetQueryBuilderFormat())

	createUc := prompt.NewCreatePromptUseCase(promptRepo)
	deleteUc := prompt.NewDeletePromptUseCase(promptRepo)
	updateUc := prompt.NewUpdatePromptUseCase(promptRepo)
	readUc := prompt.NewReadPromptUseCase(promptRepo)

	createHandler := http.NewCreateHandler(createUc)
	deleteHandler := http.NewDeleteHandler(deleteUc)
	updateHandler := http.NewUpdateHandler(updateUc)
	readHandler := http.NewReadHandler(readUc)
	readAllHandler := http.NewReadAllHandler(readUc)

	container := &admin.HandlerContainer{
		CreateHandler:  createHandler,
		ReadHandler:    readHandler,
		ReadAllHandler: readAllHandler,
		DeleteHandler:  deleteHandler,
		UpdateHandler:  updateHandler,
	}

	apiHandlers := router.NewApiHandler(container)

	go func() {
		handlers, err := apiHandlers.SetupRoutes()

		if err != nil {
			slog.Error(fmt.Sprintf("error initializing handlers: %s", err.Error()))
		}

		if err = serv.Run(os.Getenv("APP_PORT"), handlers); err != nil {
			slog.Error(fmt.Sprintf("error running server: %s", err.Error()))
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err = serv.Stop(shutdownCtx); err != nil {
		slog.ErrorContext(shutdownCtx, fmt.Sprintf("http server shutdown error: %v", err))
		return err
	}

	return nil
}
