package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Geriler/effective-mobile/internal/config"
	"github.com/Geriler/effective-mobile/internal/middleware"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPGW struct {
	cfg    config.Config
	log    *slog.Logger
	server *http.Server
	mux    *runtime.ServeMux
}

func NewHTTPGW(cfg config.Config, log *slog.Logger) *HTTPGW {
	mux := runtime.NewServeMux()

	return &HTTPGW{
		cfg: cfg,
		log: log,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
			Handler: middleware.NewLogWrapperHandler(mux, log),
		},
		mux: mux,
	}
}

func (a *HTTPGW) ListenAndServe() error {
	conn, err := grpc.NewClient(fmt.Sprintf("dns:%s:%d", a.cfg.GRPC.Host, a.cfg.GRPC.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.mux.HandlePath("GET", "/swagger", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./docs/swagger/index.html")
	})
	a.mux.HandlePath("GET", "/swagger.json", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/api/subscriptions.swagger.json")
	})

	err = pbSubscription.RegisterSubscriptionsHandler(context.Background(), a.mux, conn)
	if err != nil {
		return err
	}

	if err = a.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *HTTPGW) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
