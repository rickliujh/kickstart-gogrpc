package server

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
	grpcimpl "github.com/rickliujh/kickstart-gogrpc/pkg/api/grpc/impl"
	"github.com/rickliujh/kickstart-gogrpc/pkg/api/grpc/pb/v1/pbv1connect"
	"github.com/rickliujh/kickstart-gogrpc/pkg/sql"
	"github.com/rickliujh/kickstart-gogrpc/pkg/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func StartGRPC(addr, name, env, dbConnStr string) {

	logger := utils.NewLogger(0)

	// create server
	logger.Info("creating server...")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbConnStr)
	if err != nil {
		logger.Error(err, "unable to connect to db")
		return
	}
	defer conn.Close(ctx)
	queries := sql.New(conn)

	s, err := grpcimpl.NewServer(name, version, env, queries)
	if err != nil {
		logger.Error(err, "error while creating server")
		return
	}

	mux := http.NewServeMux()
	path, handler := pbv1connect.NewServiceHandler(s)
	mux.Handle(path, handler)

	// run server
	logger.Info("starting server...", "server_name", s.String())
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		logger.Error(err, "error while running server")
		return
	}

	logger.Info("done")
}
