package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rickliujh/kickstart-gogrpc/pkg/service"
)

func CounterAPIs(r *chi.Mux) {
	counter := service.Counter{}

	r.Route("/count", func(r chi.Router) {

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write(fmt.Append(nil, counter.Count()))
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			counter.Increment()
		})

	})
}
