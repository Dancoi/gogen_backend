package repository

import (
	"context"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionPlanRepository struct {
	queries *sqlc.Queries
}

// NewSubscriptionPlanRepository создаёт новый SubscriptionPlanRepository
func NewSubscriptionPlanRepository(pool *pgxpool.Pool) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		queries: sqlc.New(pool),
	}
}

// GetSubscriptionPlanByType получает план подписки по типу (trial, premium, commercial)
func (r *subscriptionPlanRepository) GetSubscriptionPlanByType(ctx context.Context, planType string) (*sqlc.SubscriptionPlan, error) {
	plan, err := r.queries.GetSubscriptionPlanByType(ctx, planType)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetSubscriptionPlanByID получает план подписки по ID
func (r *subscriptionPlanRepository) GetSubscriptionPlanByID(ctx context.Context, planID int32) (*sqlc.SubscriptionPlan, error) {
	plan, err := r.queries.GetSubscriptionPlanByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ListSubscriptionPlans получает все активные планы подписки
func (r *subscriptionPlanRepository) ListSubscriptionPlans(ctx context.Context) ([]sqlc.SubscriptionPlan, error) {
	return r.queries.ListSubscriptionPlans(ctx)
}
