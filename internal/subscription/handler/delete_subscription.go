package handler

import (
	"context"

	"buf.build/go/protovalidate"
	pbSubscription "github.com/Geriler/effective-mobile/pb/api"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *SubscriptionHandler) DeleteSubscription(ctx context.Context, request *pbSubscription.DeleteSubscriptionRequest) (*pbSubscription.DeleteSubscriptionResponse, error) {
	const op = "SubscriptionHandler.DeleteSubscription"
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

	err = s.service.DeleteSubscription(ctx, subscriptionID)
	if err != nil {
		logger.Error("failed delete subscription", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbSubscription.DeleteSubscriptionResponse{}, nil
}
