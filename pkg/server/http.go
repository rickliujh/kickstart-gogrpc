package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/rickliujh/kickstart-gogrpc/pkg/api/http/v1"
)

func HTTPServer(addr, name, env, dbConnStr string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write(fmt.Appendf(nil, "Server: %s, Version: %s, Environment: %s", name, version, env))
	})

	v1.Route(r)

	http.ListenAndServe(addr, r)
}
