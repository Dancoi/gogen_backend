package repository

import (
	"context"
	"time"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userSessionRepository struct {
	queries *sqlc.Queries
}

// NewUserSessionRepository создаёт новый UserSessionRepository
func NewUserSessionRepository(pool *pgxpool.Pool) UserSessionRepository {
	return &userSessionRepository{
		queries: sqlc.New(pool),
	}
}

// CreateUserSession создаёт новую сессию
func (r *userSessionRepository) CreateUserSession(ctx context.Context, userID int32, sessionToken, ipAddress, userAgent string, expiresAt int64) (*sqlc.UserSession, error) {
	sess, err := r.queries.CreateUserSession(ctx, sqlc.CreateUserSessionParams{
		UserID:       userID,
		SessionToken: sessionToken,
		IpAddress:    pgtype.Text{String: ipAddress, Valid: ipAddress != ""},
		UserAgent:    pgtype.Text{String: userAgent, Valid: userAgent != ""},
		ExpiresAt:    pgtype.Timestamp{Time: time.Unix(0, expiresAt), Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// GetUserSessionByToken получает сессию по токену
func (r *userSessionRepository) GetUserSessionByToken(ctx context.Context, sessionToken string) (*sqlc.UserSession, error) {
	sess, err := r.queries.GetUserSessionByToken(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// GetUserSessionByID получает сессию по ID
func (r *userSessionRepository) GetUserSessionByID(ctx context.Context, sessionID int32) (*sqlc.UserSession, error) {
	sess, err := r.queries.GetUserSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// ListUserSessions получает все активные сессии пользователя
func (r *userSessionRepository) ListUserSessions(ctx context.Context, userID int32) ([]sqlc.UserSession, error) {
	return r.queries.ListUserSessions(ctx, userID)
}

// UpdateUserSessionActivity обновляет время последней активности
func (r *userSessionRepository) UpdateUserSessionActivity(ctx context.Context, sessionID int32) (*sqlc.UserSession, error) {
	sess, err := r.queries.UpdateUserSessionActivity(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// RevokeUserSession отзывает конкретную сессию
func (r *userSessionRepository) RevokeUserSession(ctx context.Context, sessionID int32) error {
	return r.queries.RevokeUserSession(ctx, sessionID)
}

// RevokeAllUserSessions отзывает все сессии пользователя
func (r *userSessionRepository) RevokeAllUserSessions(ctx context.Context, userID int32) error {
	return r.queries.RevokeAllUserSessions(ctx, userID)
}

// DeleteExpiredSessions удаляет истёкшие сессии
func (r *userSessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.queries.DeleteExpiredSessions(ctx)
}
