package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Geriler/effective-mobile/internal/subscription/model"
	"github.com/Geriler/effective-mobile/pkg/lib/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get postgres connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pgxpool: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	err = goose.Up(db, "../../../migrations")
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		pool.Close()
		if err = postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgres container: %v", err)
		}
	}

	return pool, cleanup
}

func TestPostgresSubscriptionRepository_GetSumSubscriptions(t *testing.T) {
	pool, cleanup := setupTestContainer(t)
	defer cleanup()

	ctx := context.Background()
	log := logger.Setup("text", "warn")
	repo := NewPostgresSubscriptionRepository(pool, log)

	userID1, userID2 := uuid.New(), uuid.New()

	testCases := []struct {
		name          string
		subscriptions []model.Subscription
		filters       model.Filters
		expected      int32
	}{
		{
			name: "Подписка активна весь запрошенный период",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       500,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 500 * 4,
		},
		{
			name: "Подписка началась внутри периода",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       400,
					StartDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 400 * 3,
		},
		{
			name: "Подписка закончилась внутри периода",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450 * 3,
		},
		{
			name: "Подписка начинается и заканчивается внутри периода",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450 * 4,
		},
		{
			name: "Бесконечная подписка",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450 * 12,
		},
		{
			name: "Несколько подписок",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID1,
					ServiceName: "YouTube",
					Price:       300,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450*4 + 300*4,
		},
		{
			name: "Фильтр по наименованию подписки",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID1,
					ServiceName: "YouTube",
					Price:       300,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				ServiceName: "Netflix",
			},
			expected: 450 * 4,
		},
		{
			name: "Фильтр по user_id",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID2,
					ServiceName: "YouTube",
					Price:       300,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				UserID:    userID1,
			},
			expected: 450 * 4,
		},
		{
			name: "Подписки не пересекается с фильтром",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID2,
					ServiceName: "YouTube",
					Price:       300,
					StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 0,
		},
		{
			name: "Фильтр на один месяц",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450,
		},
		{
			name: "Несколько разных сервисов",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID1,
					ServiceName: "YouTube",
					Price:       300,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      userID1,
					ServiceName: "Hulu",
					Price:       200,
					StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450 + 300*3 + 200*5,
		},
		{
			name: "Бесконечная подписка, фильтр затрагивает несколько лет",
			subscriptions: []model.Subscription{
				{
					UserID:      userID1,
					ServiceName: "Netflix",
					Price:       450,
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			filters: model.Filters{
				StartDate: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: 450 * 7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, "TRUNCATE TABLE subscriptions")
			assert.NoError(t, err)

			for _, sub := range tc.subscriptions {
				_, err = repo.CreateSubscription(ctx, sub)
				assert.NoError(t, err)
			}

			sum, err := repo.GetSumSubscriptions(ctx, tc.filters)
			assert.Equal(t, tc.expected, sum)
			assert.NoError(t, err)
		})
	}
}
