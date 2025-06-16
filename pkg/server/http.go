package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	// set at build time
	version = "v0.0.1-default"
)

func HTTPServer(addr, name, env, dbConnStr string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Server: %s, Version: %s, Environment: %s", name, version, env)))
	})
	http.ListenAndServe(addr, r)
}
