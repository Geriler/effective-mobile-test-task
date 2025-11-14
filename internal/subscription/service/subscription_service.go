package service

import (
	"context"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	GetSubscriptionById(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	AllSubscriptions(ctx context.Context, pagination model.Pagination) ([]model.Subscription, error)
	UpdateSubscription(ctx context.Context, id uuid.UUID, subscription model.Subscription) (*model.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetSumSubscriptions(ctx context.Context, filters model.Filters) (int32, error)
}

type SubscriptionService struct {
	repo SubscriptionRepository
}

func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) AddSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error) {
	return s.repo.CreateSubscription(ctx, subscription)
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetSubscriptionById(ctx, id)
}

func (s *SubscriptionService) ListSubscriptions(ctx context.Context, pagination model.Pagination) ([]model.Subscription, error) {
	return s.repo.AllSubscriptions(ctx, pagination)
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, id uuid.UUID, subscription model.Subscription) (*model.Subscription, error) {
	return s.repo.UpdateSubscription(ctx, id, subscription)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSubscription(ctx, id)
}

func (s *SubscriptionService) GetTotalSum(ctx context.Context, filters model.Filters) (int32, error) {
	return s.repo.GetSumSubscriptions(ctx, filters)
}
