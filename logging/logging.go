package logging

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	jLogger := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jLogger)
	slog.SetDefault(logger)
	return logger
}
