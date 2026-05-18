package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	queries *sqlc.Queries
}

// NewSubscriptionRepository создаёт новый SubscriptionRepository
func NewSubscriptionRepository(pool *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{
		queries: sqlc.New(pool),
	}
}

// CreateSubscription создаёт новую подписку
func (r *subscriptionRepository) CreateSubscription(ctx context.Context, userID, planID int32, planType string, expiresAt int64, maxTokens int32) (*sqlc.Subscription, error) {
	fmt.Printf("[CreateSubscription] Creating subscription for user: %d, plan: %d, expiresAt: %v\n", userID, planID, time.Unix(expiresAt, 0))

	sub, err := r.queries.CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
		UserID:             userID,
		SubscriptionPlanID: planID,
		PlanType:           planType,
		IsActive:           pgtype.Bool{Bool: true, Valid: true},
		ExpiresAt:          pgtype.Timestamp{Time: time.Unix(expiresAt, 0), Valid: true},
		MaxTokensPerMonth:  maxTokens,
	})
	if err != nil {
		fmt.Printf("[CreateSubscription] Error creating subscription: %v\n", err)
		return nil, err
	}
	fmt.Printf("[CreateSubscription] Subscription created successfully: ID=%s\n", sub.ID)
	return &sub, nil
}

// GetUserActiveSubscription получает активную подписку пользователя
func (r *subscriptionRepository) GetUserActiveSubscription(ctx context.Context, userID int32) (*sqlc.Subscription, error) {
	sub, err := r.queries.GetUserActiveSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetUserSubscription получает последнюю подписку пользователя
func (r *subscriptionRepository) GetUserSubscription(ctx context.Context, userID int32) (*sqlc.Subscription, error) {
	sub, err := r.queries.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetSubscriptionByID получает подписку по ID
func (r *subscriptionRepository) GetSubscriptionByID(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error) {
	sub, err := r.queries.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// IncrementSubscriptionUsage увеличивает счётчик использования токенов
func (r *subscriptionRepository) IncrementSubscriptionUsage(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error) {
	sub, err := r.queries.IncrementSubscriptionUsage(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// ResetSubscriptionUsage сбрасывает счётчик использования токенов
func (r *subscriptionRepository) ResetSubscriptionUsage(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error) {
	sub, err := r.queries.ResetSubscriptionUsage(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// UpdateSubscriptionIsActive обновляет статус активности подписки
func (r *subscriptionRepository) UpdateSubscriptionIsActive(ctx context.Context, subscriptionID int32, isActive bool) (*sqlc.Subscription, error) {
	sub, err := r.queries.UpdateSubscriptionIsActive(ctx, sqlc.UpdateSubscriptionIsActiveParams{
		ID:       subscriptionID,
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
