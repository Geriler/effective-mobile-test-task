package handler

import (
	"context"
	"errors"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SubscriptionHandler) GetSubscriptions(ctx context.Context, request *pbSubscription.GetSubscriptionsRequest) (*pbSubscription.GetSubscriptionsResponse, error) {
	const op = "SubscriptionHandler.GetSubscriptions"
	logger := s.logger.With("op", op).With("request", request)

	subscriptions, err := s.service.ListSubscriptions(ctx)
	if err != nil {
		if errors.Is(err, model.ErrSubscriptionNotFound) {
			return &pbSubscription.GetSubscriptionsResponse{
				Subscriptions: []*pbSubscription.Subscription{},
			}, nil
		}

		logger.Error("failed to list subscriptions", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	subscriptionsResponse := make([]*pbSubscription.Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		subscriptionResponse := pbSubscription.Subscription{
			SubscriptionId: subscription.ID.String(),
			UserId:         subscription.UserID.String(),
			ServiceName:    subscription.ServiceName,
			Price:          subscription.Price,
			StartDate:      subscription.StartDate.Format("01-2006"),
		}

		if !subscription.EndDate.IsZero() {
			endDate := subscription.EndDate.Format("01-2006")
			subscriptionResponse.EndDate = &endDate
		}

		subscriptionsResponse = append(subscriptionsResponse, &subscriptionResponse)
	}

	return &pbSubscription.GetSubscriptionsResponse{
		Subscriptions: subscriptionsResponse,
	}, nil
}
