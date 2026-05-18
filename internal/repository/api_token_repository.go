package repository

import (
	"context"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type apiTokenRepository struct {
	queries *sqlc.Queries
}

// NewApiTokenRepository создаёт новый ApiTokenRepository
func NewApiTokenRepository(pool *pgxpool.Pool) ApiTokenRepository {
	return &apiTokenRepository{
		queries: sqlc.New(pool),
	}
}

// CreateApiToken создаёт новый API токен
func (r *apiTokenRepository) CreateApiToken(ctx context.Context, userID, subscriptionID int32, name, tokenHash string) (*sqlc.ApiToken, error) {
	token, err := r.queries.CreateApiToken(ctx, sqlc.CreateApiTokenParams{
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Name:           name,
		TokenHash:      tokenHash,
		IsActive:       pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetApiTokenByHash получает токен по хешу
func (r *apiTokenRepository) GetApiTokenByHash(ctx context.Context, tokenHash string) (*sqlc.ApiToken, error) {
	token, err := r.queries.GetApiTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetApiTokenByID получает токен по ID
func (r *apiTokenRepository) GetApiTokenByID(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error) {
	token, err := r.queries.GetApiTokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// ListUserApiTokens получает все токены пользователя
func (r *apiTokenRepository) ListUserApiTokens(ctx context.Context, userID int32) ([]sqlc.ApiToken, error) {
	return r.queries.ListUserApiTokens(ctx, userID)
}

// UpdateApiTokenLastUsed обновляет время последнего использования токена
func (r *apiTokenRepository) UpdateApiTokenLastUsed(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error) {
	token, err := r.queries.UpdateApiTokenLastUsed(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// RevokeApiToken отзывает токен (mark as revoked)
func (r *apiTokenRepository) RevokeApiToken(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error) {
	token, err := r.queries.RevokeApiToken(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// CountUserActiveTokens считает активные токены пользователя
func (r *apiTokenRepository) CountUserActiveTokens(ctx context.Context, userID int32) (int64, error) {
	return r.queries.CountUserActiveTokens(ctx, userID)
}

// DeleteApiToken удаляет токен
func (r *apiTokenRepository) DeleteApiToken(ctx context.Context, tokenID int32) error {
	return r.queries.DeleteApiToken(ctx, tokenID)
}
