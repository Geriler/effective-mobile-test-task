package handler

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/Geriler/effective-mobile/internal/subscription/model"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SubscriptionHandler) UpdateSubscription(ctx context.Context, request *pbSubscription.UpdateSubscriptionRequest) (*pbSubscription.UpdateSubscriptionResponse, error) {
	const op = "SubscriptionHandler.UpdateSubscription"
	logger := s.logger.With("op", op).With("request", request)

	err := protovalidate.Validate(request)
	if err != nil {
		logger.Error("failed validate request", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	subscriptionID, err := uuid.Parse(request.GetSubscriptionId())
	if err != nil {
		logger.Error("failed parse subscription id", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var endDate time.Time
	if request.GetEndDate() != "" {
		endDate, err = time.Parse("01-2006", request.GetEndDate())
		if err != nil {
			logger.Error("failed parse end date", "error", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	subscription := model.Subscription{
		ServiceName: request.GetServiceName(),
		Price:       request.GetPrice(),
		EndDate:     endDate,
	}

	resultSubscription, err := s.service.UpdateSubscription(ctx, subscriptionID, subscription)
	if err != nil {
		logger.Error("failed update subscription", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pbSubscription.UpdateSubscriptionResponse{
		Subscription: &pbSubscription.Subscription{
			SubscriptionId: resultSubscription.ID.String(),
			UserId:         resultSubscription.UserID.String(),
			ServiceName:    resultSubscription.ServiceName,
			Price:          resultSubscription.Price,
			StartDate:      resultSubscription.StartDate.Format("01-2006"),
		},
	}

	if !subscription.EndDate.IsZero() {
		endDate := subscription.EndDate.Format("01-2006")
		response.Subscription.EndDate = &endDate
	}

	return response, nil
}
