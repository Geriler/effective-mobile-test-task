package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Geriler/effective-mobile/internal/app"
	"github.com/Geriler/effective-mobile/internal/config"
	"github.com/Geriler/effective-mobile/pkg/lib/logger"
)

func main() {
	rootCtx := context.Background()

	cfg := config.MustLoad()

	log := logger.Setup(cfg.Logger.Type, cfg.Logger.Level)

	grpcServer, err := app.NewGRPCServer(rootCtx, cfg, log)
	if err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	httpgwServer := app.NewHTTPGW(cfg, log)

	go func() {
		log.Info("starting gRPC application", "port", cfg.GRPC.Port)
		err = grpcServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("starting HTTP application", "port", cfg.HTTP.Port)
		err = httpgwServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(rootCtx, cfg.TimeoutStop)
	defer cancel()

	grpcServer.Shutdown()
	err = httpgwServer.Shutdown(ctx)
	if err != nil {
		log.Error(err.Error())
	}
}
