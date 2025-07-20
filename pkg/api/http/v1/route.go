package v1

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func Route(r *chi.Mux, logger *slog.Logger) {
	CounterAPIs(r, logger)
}
