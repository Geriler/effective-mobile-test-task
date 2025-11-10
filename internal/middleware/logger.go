package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Geriler/effective-mobile/pkg/lib/logger"
	"google.golang.org/grpc"
)

type LogWrapperHandler struct {
	wrap   http.Handler
	logger *slog.Logger
}

func NewLogWrapperHandler(wrap http.Handler, logger *slog.Logger) *LogWrapperHandler {
	return &LogWrapperHandler{
		wrap:   wrap,
		logger: logger,
	}
}

func (h LogWrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "middleware.LogWrapperHandler.ServeHTTP"

	log := h.logger.With(
		slog.String("op", op),
		slog.String("url", r.URL.String()),
		slog.String("method", r.Method),
	)

	log.Info("request received")

	h.wrap.ServeHTTP(w, r)
}

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	const op = "middleware.Logger"

	log := logger.GetLogger().With(
		slog.String("op", op),
		slog.String("method", info.FullMethod),
	)

	log.Info("request received")

	resp, err := handler(ctx, req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp, err
}
