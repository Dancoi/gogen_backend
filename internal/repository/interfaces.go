package repository

import (
	"context"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, email, username, passwordHash string) (*sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error)
	GetUserByID(ctx context.Context, id int32) (*sqlc.User, error)
	GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error)
	UpdateUserPasswordHash(ctx context.Context, userID int32, passwordHash string) (*sqlc.User, error)
	UpdateUserIsActive(ctx context.Context, userID int32, isActive bool) (*sqlc.User, error)
	UpdateUserEmailVerified(ctx context.Context, userID int32, emailVerified bool) (*sqlc.User, error)
	DeleteUser(ctx context.Context, userID int32) error
}

// SubscriptionRepository интерфейс для работы с подписками
type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, userID, planID int32, planType string, expiresAt int64, maxTokens int32) (*sqlc.Subscription, error)
	GetUserActiveSubscription(ctx context.Context, userID int32) (*sqlc.Subscription, error)
	GetUserSubscription(ctx context.Context, userID int32) (*sqlc.Subscription, error)
	GetSubscriptionByID(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error)
	IncrementSubscriptionUsage(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error)
	ResetSubscriptionUsage(ctx context.Context, subscriptionID int32) (*sqlc.Subscription, error)
	UpdateSubscriptionIsActive(ctx context.Context, subscriptionID int32, isActive bool) (*sqlc.Subscription, error)
}

// ApiTokenRepository интерфейс для работы с API токенами
type ApiTokenRepository interface {
	CreateApiToken(ctx context.Context, userID, subscriptionID int32, name, tokenHash string) (*sqlc.ApiToken, error)
	GetApiTokenByHash(ctx context.Context, tokenHash string) (*sqlc.ApiToken, error)
	GetApiTokenByID(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error)
	ListUserApiTokens(ctx context.Context, userID int32) ([]sqlc.ApiToken, error)
	UpdateApiTokenLastUsed(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error)
	RevokeApiToken(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error)
	CountUserActiveTokens(ctx context.Context, userID int32) (int64, error)
	DeleteApiToken(ctx context.Context, tokenID int32) error
}

// UserSessionRepository интерфейс для работы с сессиями
type UserSessionRepository interface {
	CreateUserSession(ctx context.Context, userID int32, sessionToken, ipAddress, userAgent string, expiresAt int64) (*sqlc.UserSession, error)
	GetUserSessionByToken(ctx context.Context, sessionToken string) (*sqlc.UserSession, error)
	GetUserSessionByID(ctx context.Context, sessionID int32) (*sqlc.UserSession, error)
	ListUserSessions(ctx context.Context, userID int32) ([]sqlc.UserSession, error)
	UpdateUserSessionActivity(ctx context.Context, sessionID int32) (*sqlc.UserSession, error)
	RevokeUserSession(ctx context.Context, sessionID int32) error
	RevokeAllUserSessions(ctx context.Context, userID int32) error
	DeleteExpiredSessions(ctx context.Context) error
}

// AuditLogRepository интерфейс для работы с логами аудита
type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, userID int32, action, resourceType string, resourceID int32, changes []byte, statusCode int32, errorMessage, ipAddress, userAgent string) (*sqlc.AuditLog, error)
	GetUserAuditLogs(ctx context.Context, userID int32, limit int32, offset int32) ([]sqlc.AuditLog, error)
}

// SubscriptionPlanRepository интерфейс для работы с планами подписки
type SubscriptionPlanRepository interface {
	GetSubscriptionPlanByType(ctx context.Context, planType string) (*sqlc.SubscriptionPlan, error)
	GetSubscriptionPlanByID(ctx context.Context, planID int32) (*sqlc.SubscriptionPlan, error)
	ListSubscriptionPlans(ctx context.Context) ([]sqlc.SubscriptionPlan, error)
}
