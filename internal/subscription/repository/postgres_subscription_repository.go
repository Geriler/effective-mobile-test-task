package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	repository "github.com/Geriler/effective-mobile/internal/subscription/repository/sqlc"
	"github.com/Geriler/effective-mobile/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSubscriptionRepository struct {
	conn   *pgxpool.Pool
	cmd    *repository.Queries
	logger *slog.Logger
}

func NewPostgresSubscriptionRepository(conn *pgxpool.Pool, logger *slog.Logger) *PostgresSubscriptionRepository {
	cmd := repository.New(conn)

	return &PostgresSubscriptionRepository{
		conn:   conn,
		cmd:    cmd,
		logger: logger,
	}
}

func (r *PostgresSubscriptionRepository) CreateSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error) {
	// На будущее - можно добавить транзакции, если потребуется, например, outbox
	const op = "PostgresSubscriptionRepository.CreateSubscription"
	logger := r.logger.With("op", op).With("subscription", subscription)

	row, err := r.cmd.CreateSubscription(ctx, repository.CreateSubscriptionParams{
		UserID:      utils.GoogleUUIDToPgxUUID(subscription.UserID),
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		StartDate: pgtype.Date{
			Time:  subscription.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  subscription.EndDate,
			Valid: !subscription.EndDate.IsZero(),
		},
	})
	if err != nil {
		logger.Error("failed to create subscription", "error", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, model.ErrSubscriptionAlreadyExists
		}
		return nil, err
	}

	subscriptionID, err := utils.PgxUUIDToGoogleUUID(row.ID)
	if err != nil {
		logger.Warn("failed to convert pgx UUID to google UUID (subscription_id)", "error", err)
		return nil, err
	}

	userID, err := utils.PgxUUIDToGoogleUUID(row.UserID)
	if err != nil {
		logger.Warn("failed to convert pgx UUID to google UUID (user_id)", "error", err)
		return nil, err
	}

	return &model.Subscription{
		ID:          subscriptionID,
		UserID:      userID,
		ServiceName: row.ServiceName,
		Price:       row.Price,
		StartDate:   row.StartDate.Time,
		EndDate:     row.EndDate.Time,
	}, nil
}

func (r *PostgresSubscriptionRepository) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	const op = "PostgresSubscriptionRepository.GetSubscriptionById"
	logger := r.logger.With("op", op).With("subscription_id", id)

	subscription, err := r.cmd.GetSubscriptionById(ctx, utils.GoogleUUIDToPgxUUID(id))
	if errors.Is(err, sql.ErrNoRows) {
		logger.Warn("subscription not found")
		return nil, model.ErrSubscriptionNotFound
	}
	if err != nil {
		logger.Error("failed to get subscription", "error", err)
		return nil, err
	}

	userID, err := utils.PgxUUIDToGoogleUUID(subscription.UserID)
	if err != nil {
		logger.Error("failed to convert pgx UUID to google UUID", "error", err)
		return nil, err
	}

	return &model.Subscription{
		ID:          id,
		UserID:      userID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		StartDate:   subscription.StartDate.Time,
		EndDate:     subscription.EndDate.Time,
	}, nil
}

func (r *PostgresSubscriptionRepository) AllSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	const op = "PostgresSubscriptionRepository.AllSubscriptions"
	logger := r.logger.With("op", op)

	rows, err := r.cmd.AllSubscriptions(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Warn("subscriptions not found")
		return nil, model.ErrSubscriptionNotFound
	}
	if err != nil {
		logger.Error("failed to get all subscriptions", "error", err)
		return nil, err
	}

	subscriptions := make([]model.Subscription, 0, len(rows))
	for _, row := range rows {
		subscriptionID, err := utils.PgxUUIDToGoogleUUID(row.ID)
		if err != nil {
			logger.Error("failed to convert pgx UUID to google UUID (subscription_id)", "error", err)
			return nil, err
		}

		userID, err := utils.PgxUUIDToGoogleUUID(row.UserID)
		if err != nil {
			logger.Error("failed to convert pgx UUID to google UUID (user_id)", "error", err)
			return nil, err
		}

		subscription := model.Subscription{
			ID:          subscriptionID,
			UserID:      userID,
			ServiceName: row.ServiceName,
			Price:       row.Price,
			StartDate:   row.StartDate.Time,
			EndDate:     row.EndDate.Time,
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r *PostgresSubscriptionRepository) UpdateSubscription(ctx context.Context, id uuid.UUID, subscription model.Subscription) (*model.Subscription, error) {
	// На будущее - можно добавить транзакции, если потребуется, например, outbox
	const op = "PostgresSubscriptionRepository.UpdateSubscription"
	logger := r.logger.With("op", op).With("subscription_id", id)

	params := repository.UpdateSubscriptionParams{
		SubscriptionID: utils.GoogleUUIDToPgxUUID(id),
		ServiceName: pgtype.Text{
			String: subscription.ServiceName,
			Valid:  subscription.ServiceName != "",
		},
		Price: pgtype.Int4{
			Int32: subscription.Price,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  subscription.EndDate,
			Valid: !subscription.EndDate.IsZero(),
		},
	}

	row, err := r.cmd.UpdateSubscription(ctx, params)
	if err != nil {
		logger.Error("failed to update subscription", "error", err)

		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr):
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, model.ErrSubscriptionAlreadyExists
			}
		case errors.Is(err, pgx.ErrNoRows):
			return nil, model.ErrSubscriptionNotFound
		default:
			return nil, err
		}
	}

	subscriptionID, err := utils.PgxUUIDToGoogleUUID(row.ID)
	if err != nil {
		logger.Warn("failed to convert pgx UUID to google UUID (subscription_id)", "error", err)
		return nil, err
	}

	userID, err := utils.PgxUUIDToGoogleUUID(row.UserID)
	if err != nil {
		logger.Warn("failed to convert pgx UUID to google UUID (user_id)", "error", err)
		return nil, err
	}

	return &model.Subscription{
		ID:          subscriptionID,
		UserID:      userID,
		ServiceName: row.ServiceName,
		Price:       row.Price,
		StartDate:   row.StartDate.Time,
		EndDate:     row.EndDate.Time,
	}, nil
}

func (r *PostgresSubscriptionRepository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	// На будущее - можно добавить транзакции, если потребуется, например, outbox
	const op = "PostgresSubscriptionRepository.DeleteSubscription"
	logger := r.logger.With("op", op).With("subscription_id", id)

	err := r.cmd.DeleteSubscription(ctx, utils.GoogleUUIDToPgxUUID(id))
	if err != nil {
		logger.Error("failed to delete subscription", "error", err)
		return err
	}

	return nil
}

func (r *PostgresSubscriptionRepository) GetSumSubscriptions(ctx context.Context, filters model.Filters) (int32, error) {
	const op = "PostgresSubscriptionRepository.GetSubscriptionsByFilters"
	logger := r.logger.With("op", op).With("filters", filters)

	params := repository.GetSumSubscriptionsParams{
		StartDate: pgtype.Date{
			Time:  filters.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  filters.EndDate,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: filters.UserID,
			Valid: filters.UserID != uuid.Nil,
		},
		ServiceName: pgtype.Text{
			String: filters.ServiceName,
			Valid:  filters.ServiceName != "",
		},
	}

	sum, err := r.cmd.GetSumSubscriptions(ctx, params)
	if err != nil {
		logger.Error("failed to get sum subscriptions", "error", err)
		return 0, err
	}

	return sum, nil
}
