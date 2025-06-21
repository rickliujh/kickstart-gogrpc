package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
	v1 "github.com/rickliujh/kickstart-gogrpc/pkg/api/http/v1"
	"github.com/rickliujh/kickstart-gogrpc/pkg/utils"
)

func HTTPServer(addr, name, env, dbConnStr, levelStr string, isLocalhost, isDebugHeaderSet bool) {
	r := chi.NewRouter()

	level, err := utils.ParseSlogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	// configure logging middleware
	logFormat := httplog.SchemaECS.Concise(isLocalhost)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
		Level:       level,
	})).With(
		slog.String("app", name),
		slog.String("version", version),
		slog.String("env", env),
	)

	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		Level:         level,
		Schema:        httplog.SchemaECS,
		RecoverPanics: false,
		Skip: func(req *http.Request, respStatus int) bool {
			return respStatus == 404 || respStatus == 405
		},
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},
		LogRequestBody:     func(req *http.Request) bool { return isDebugHeaderSet },
		LogResponseBody:    func(req *http.Request) bool { return isDebugHeaderSet },
	}))

	// heath check
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write(fmt.Appendf(nil, "Server: %s, Version: %s, Environment: %s", name, version, env))
	})

	v1.Route(r)

	http.ListenAndServe(addr, r)
}
