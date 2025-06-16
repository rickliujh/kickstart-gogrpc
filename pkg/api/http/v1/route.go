package v1

import "github.com/go-chi/chi/v5"

func Route(r *chi.Mux) {
	CounterAPIs(r)
}
