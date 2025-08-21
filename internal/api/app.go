package api

import (
	"context"
	"fmt"
	"github.com/Base-111/backend/internal/api/router"
	"github.com/Base-111/backend/internal/api/server"
	"github.com/Base-111/backend/internal/config"
	"github.com/Base-111/backend/internal/entities/admin/controller/http"
	admin "github.com/Base-111/backend/internal/entities/admin/controller/http/router"
	"github.com/Base-111/backend/internal/entities/admin/repository/postgres"
	"github.com/Base-111/backend/internal/entities/admin/usecase/prompt"
	authhandler "github.com/Base-111/backend/internal/entities/auth/controller/http"
	rds2 "github.com/Base-111/backend/internal/entities/auth/store/redis"
	"github.com/Base-111/backend/pkg/auth"
	"github.com/Base-111/backend/pkg/repository"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message/translations/english"
	"github.com/muonsoft/validation/validator"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	ctx := context.Background()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}

	serverAddr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	err = validator.SetUp(
		validation.DefaultLanguage(language.English),
		validation.Translations(english.Messages),
	)

	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&cfg.RedisConfig)
	if err = redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
		return err
	}
	defer redisClient.Close()

	serv := new(server.Server)

	authOptions := []auth.Option{
		auth.WithClientSecret(cfg.Auth.ClientSecret),
		auth.WithRealmKeycloak(cfg.Auth.Realm),
	}
	authClient, err := auth.New(
		ctx,
		cfg.Auth.BaseURL,
		cfg.Auth.ClientID,
		cfg.Auth.RedirectURL,
		authOptions...,
	)

	if err != nil {
		return err
	}

	authStore := rds2.NewAuthRedisManager(redisClient)
	sessionStore := rds2.NewSessionRedisManager(redisClient)
	authHandler := authhandler.New(cfg,
		serverAddr,
		authClient,
		authStore,
		sessionStore,
	)

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

	apiHandlers := router.NewApiHandler(authHandler, container)

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
