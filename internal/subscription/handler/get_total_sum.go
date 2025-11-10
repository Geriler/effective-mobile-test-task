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

func (s *SubscriptionHandler) GetSumSubscriptions(ctx context.Context, request *pbSubscription.GetSumSubscriptionsRequest) (*pbSubscription.GetSumSubscriptionsResponse, error) {
	const op = "SubscriptionHandler.GetSumSubscriptions"
	logger := s.logger.With("op", op).With("request", request)

	err := protovalidate.Validate(request)
	if err != nil {
		logger.Error("failed validate request", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	startDate, err := time.Parse("01-2006", request.GetStartDate())
	if err != nil {
		logger.Error("failed parse start date", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	endDate, err := time.Parse("01-2006", request.GetEndDate())
	if err != nil {
		logger.Error("failed parse start date", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var userID uuid.UUID
	if request.GetUserId() != "" {
		userID, err = uuid.Parse(request.GetUserId())
		if err != nil {
			logger.Error("failed parse user id", "error", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	sum, err := s.service.GetTotalSum(ctx, model.Filters{
		StartDate:   startDate,
		EndDate:     endDate,
		UserID:      userID,
		ServiceName: request.GetServiceName(),
	})
	if err != nil {
		logger.Error("failed get total sum", "error", err)
		return nil, err
	}

	return &pbSubscription.GetSumSubscriptionsResponse{
		TotalSum: sum,
	}, nil
}
