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

func (s *SubscriptionHandler) AddSubscription(ctx context.Context, request *pbSubscription.AddSubscriptionRequest) (*pbSubscription.AddSubscriptionResponse, error) {
	const op = "SubscriptionHandler.AddSubscription"
	logger := s.logger.With("op", op).With("request", request)

	err := protovalidate.Validate(request)
	if err != nil {
		logger.Error("failed validate request", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := uuid.Parse(request.GetUserId())
	if err != nil {
		logger.Error("failed parse user id", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	startDate, err := time.Parse("01-2006", request.GetStartDate())
	if err != nil {
		logger.Error("failed parse start date", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	subscription := model.Subscription{
		UserID:      userID,
		ServiceName: request.GetServiceName(),
		Price:       request.GetPrice(),
		StartDate:   startDate,
	}

	if request.GetEndDate() != "" {
		endDate, err := time.Parse("01-2006", request.GetEndDate())
		if err != nil {
			logger.Error("failed parse end date", "error", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		subscription.EndDate = endDate
	}

	resultSubscription, err := s.service.AddSubscription(ctx, subscription)
	if err != nil {
		logger.Error("failed add subscription", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pbSubscription.AddSubscriptionResponse{
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
