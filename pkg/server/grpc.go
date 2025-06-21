package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	grpcimpl "github.com/rickliujh/kickstart-gogrpc/pkg/api/grpc/impl"
	"github.com/rickliujh/kickstart-gogrpc/pkg/api/grpc/pb/v1/pbv1connect"
	"github.com/rickliujh/kickstart-gogrpc/pkg/sql"
	"github.com/rickliujh/kickstart-gogrpc/pkg/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func StartGRPC(addr, name, env, dbConnStr, levelStr string) {
	level, err := utils.ParseSlogLevel(levelStr)
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})).With(
		slog.String("app", name),
		slog.String("version", version),
		slog.String("env", env),
	)

	// create server
	logger.Info("creating server...")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbConnStr)
	if err != nil {
		logger.Error("unable to connect to db", slog.Any("error", err))
		return
	}
	defer conn.Close(ctx)
	queries := sql.New(conn)

	s, err := grpcimpl.NewServer(name, version, env, queries, logger)
	if err != nil {
		logger.Error("error while creating server", slog.Any("error", err))
		return
	}

	mux := http.NewServeMux()
	path, handler := pbv1connect.NewServiceHandler(s)
	mux.Handle(path, handler)

	// run server
	logger.Info("starting server...", slog.String("server_name", s.String()))
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		logger.Error("error while running server", slog.Any("error", err))
		return
	}
}
