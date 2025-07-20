package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rickliujh/kickstart-gogrpc/pkg/service"
)

func CounterAPIs(r *chi.Mux, logger *slog.Logger) {
	counter := service.Counter{}

	r.Route("/count", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write(fmt.Append(nil, counter.Count()))
			if err != nil {
				logger.Error("Write to response error", "err", err)
			}
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			counter.Increment()
		})
	})
}
