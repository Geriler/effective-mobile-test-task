package handler

import (
	"context"
	"errors"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SubscriptionHandler) GetSubscription(ctx context.Context, request *pbSubscription.GetSubscriptionRequest) (*pbSubscription.GetSubscriptionResponse, error) {
	const op = "SubscriptionHandler.GetSubscription"
	logger := s.logger.With("op", op).With("request", request)

	subscriptionID, err := uuid.Parse(request.GetSubscriptionId())
	if err != nil {
		logger.Error("failed to parse subscription id", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	subscription, err := s.service.GetSubscription(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, model.ErrSubscriptionNotFound) {
			return &pbSubscription.GetSubscriptionResponse{
				Subscription: nil,
			}, status.Error(codes.NotFound, err.Error())
		}

		logger.Error("failed to get subscription", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pbSubscription.GetSubscriptionResponse{
		Subscription: &pbSubscription.Subscription{
			SubscriptionId: subscription.ID.String(),
			UserId:         subscription.UserID.String(),
			ServiceName:    subscription.ServiceName,
			Price:          subscription.Price,
			StartDate:      subscription.StartDate.Format("01-2006"),
		},
	}

	if !subscription.EndDate.IsZero() {
		endDate := subscription.EndDate.Format("01-2006")
		response.Subscription.EndDate = &endDate
	}

	return response, nil
}
