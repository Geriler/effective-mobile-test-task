package handler

import (
	"context"
	"log/slog"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	AddSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context, pagination model.Pagination) ([]model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, subscription model.Subscription) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetTotalSum(ctx context.Context, filters model.Filters) (int32, error)
}

type SubscriptionHandler struct {
	pbSubscription.UnimplementedSubscriptionsServer
	logger  *slog.Logger
	service SubscriptionService
}

func NewSubscriptionHandler(logger *slog.Logger, subscriptionService SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger:  logger,
		service: subscriptionService,
	}
}
