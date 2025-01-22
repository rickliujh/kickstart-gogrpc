package utils

import (
	"log/slog"
	"os"

	"github.com/go-logr/logr"
)

func NewLogger() *logr.Logger {
	logger := logr.FromSlogHandler(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		},
	))
	return &logger
}
