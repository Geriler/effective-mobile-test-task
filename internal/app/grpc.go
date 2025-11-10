package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/Geriler/effective-mobile/internal/config"
	"github.com/Geriler/effective-mobile/internal/middleware"
	"github.com/Geriler/effective-mobile/internal/subscription/handler"
	"github.com/Geriler/effective-mobile/internal/subscription/repository"
	"github.com/Geriler/effective-mobile/internal/subscription/service"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/Geriler/effective-mobile/pkg/infra/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	cfg    config.Config
	log    *slog.Logger
	server *grpc.Server
}

func NewGRPCServer(ctx context.Context, cfg config.Config, log *slog.Logger) (*GRPCServer, error) {
	conn, err := postgres.Connect(ctx, fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name),
	)
	if err != nil {
		return nil, err
	}

	subscriptionRepo := repository.NewPostgresSubscriptionRepository(conn, log)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(log, subscriptionService)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Logger,
		),
	)

	pbSubscription.RegisterSubscriptionsServer(server, subscriptionHandler)
	reflection.Register(server)

	return &GRPCServer{
		cfg:    cfg,
		log:    log,
		server: server,
	}, nil
}

func (a *GRPCServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port))
	if err != nil {
		return err
	}

	if err = a.server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (a *GRPCServer) Shutdown() {
	a.server.GracefulStop()
}
