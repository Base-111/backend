package logs

import (
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"log/slog"
	"os"
)

func Init() error {
	gelfAddr := os.Getenv("GELF_ADDR")
	if gelfAddr == "" {
		gelfAddr = "graylog:5044"
	}
	writer, err := gelf.NewTCPWriter(gelfAddr)
	if err != nil {
		return err
	}
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler)

	slog.SetDefault(logger)
	slog.Info("Hello GrayLog")
	return nil
}
