package main

import (
	"flag"
	"net/http"

	"github.com/rickliujh/kickstart-gogrpc/pkg/api/v1/apiv1connect"
	"github.com/rickliujh/kickstart-gogrpc/pkg/server"
	"github.com/rickliujh/kickstart-gogrpc/pkg/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	address     string
	name        = "server"
	environment = "development"

	// set at build time
	version = "v0.0.1-default"
)

func main() {
	flag.StringVar(&address, "address", ":8080", "Server address (host:port)")
	flag.StringVar(&name, "name", name, "Server name (default: server)")
	flag.StringVar(&environment, "environment", environment, "Server environment (default: development)")
	flag.Parse()

	logger := utils.NewLogger()

	// create server
	logger.Info("creating server...")
	s, err := server.NewServer(name, version, environment)
	if err != nil {
		logger.Error(err, "error while creating server")
	}

	mux := http.NewServeMux()
	path, handler := apiv1connect.NewServiceHandler(s)
	mux.Handle(path, handler)

	// run server
	logger.Info("starting server: %s", s.String())
	if err := http.ListenAndServe(
		address,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		logger.Error(err, "error while running server")
	}

	logger.Info("done")
}
