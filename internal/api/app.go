package api

import (
	"context"
	"fmt"
	"github.com/Base-111/backend/internal/api/router"
	"github.com/Base-111/backend/internal/api/server"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/message/translations/english"
	"github.com/muonsoft/validation/validator"
	"golang.org/x/text/language"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {

	//err := logs.Init()
	//if err != nil {
	//	slog.Error(fmt.Sprintf("error validation starting: %s", err.Error()))
	//	return err
	//}

	err := validator.SetUp(
		validation.DefaultLanguage(language.English),
		validation.Translations(english.Messages),
	)

	if err != nil {
		return err
	}

	serv := new(server.Server)

	apiHandlers := router.NewApiHandler()

	go func() {
		handlers, err := apiHandlers.SetupRoutes()

		if err != nil {
			slog.Error(fmt.Sprintf("error initializing handlers: %s", err.Error()))
		}

		if err = serv.Run(os.Getenv("API_PORT"), handlers); err != nil {
			slog.Error(fmt.Sprintf("error running server: %s", err.Error()))
		}
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	<-osSignal

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := serv.Stop(shutdownCtx); err != nil {
		slog.ErrorContext(shutdownCtx, fmt.Sprintf("http server shutdown error: %v", err))
		return err
	}

	return nil
}
