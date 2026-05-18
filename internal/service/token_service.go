package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dancoi/gogen_backend/internal/repository"
	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/Dancoi/gogen_backend/pkg/utils/errors"
	"github.com/Dancoi/gogen_backend/pkg/utils/hash"
	"github.com/Dancoi/gogen_backend/pkg/utils/randstring"
)

// TokenService интерфейс для работы с API токенами
type TokenService interface {
	GenerateApiToken(ctx context.Context, userID int32, tokenName string) (string, error)
	ValidateApiToken(ctx context.Context, token string) (*sqlc.ApiToken, error)
	RevokeApiToken(ctx context.Context, tokenID int32) error
	ListUserTokens(ctx context.Context, userID int32) ([]sqlc.ApiToken, error)
	GetTokenByID(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error)
}

type tokenService struct {
	apiTokenRepo     repository.ApiTokenRepository
	subscriptionRepo repository.SubscriptionRepository
	auditLogRepo     repository.AuditLogRepository
}

// NewTokenService создаёт новый TokenService
func NewTokenService(
	apiTokenRepo repository.ApiTokenRepository,
	subscriptionRepo repository.SubscriptionRepository,
	auditLogRepo repository.AuditLogRepository,
) TokenService {
	return &tokenService{
		apiTokenRepo:     apiTokenRepo,
		subscriptionRepo: subscriptionRepo,
		auditLogRepo:     auditLogRepo,
	}
}

// GenerateApiToken генерирует новый API токен для пользователя
func (s *tokenService) GenerateApiToken(ctx context.Context, userID int32, tokenName string) (string, error) {
	fmt.Printf("[GenerateApiToken] Starting token generation for user: %d\n", userID)

	// Получаем активную подписку пользователя
	subscription, err := s.subscriptionRepo.GetUserActiveSubscription(ctx, userID)
	if err != nil {
		fmt.Printf("[GenerateApiToken] Error getting subscription: %v\n", err)
		if err.Error() == "no rows in result set" {
			fmt.Printf("[GenerateApiToken] No active subscription found for user: %d\n", userID)
			return "", errors.ErrSubscriptionNotFound
		}
		return "", fmt.Errorf("failed to get subscription: %w", err)
	}
	fmt.Printf("[GenerateApiToken] Subscription found: ID=%s, ExpiresAt=%v\n", subscription.ID, subscription.ExpiresAt)

	// Проверяем, не истекла ли подписка
	// SQLC использует pgtype для времени, нужно проверить через Valid и Time
	if subscription.ExpiresAt.Valid && subscription.ExpiresAt.Time.Before(time.Now()) {
		fmt.Println("[GenerateApiToken] Subscription expired")
		return "", errors.ErrSubscriptionExpired
	}
	fmt.Println("[GenerateApiToken] Subscription is active")

	// Проверяем лимит токенов
	activeTokenCount, err := s.apiTokenRepo.CountUserActiveTokens(ctx, userID)
	if err != nil {
		fmt.Printf("[GenerateApiToken] Error counting tokens: %v\n", err)
		return "", fmt.Errorf("failed to count active tokens: %w", err)
	}
	fmt.Printf("[GenerateApiToken] Active token count: %d\n", activeTokenCount)

	// Допустим максимум 10 активных токенов на пользователя
	if activeTokenCount >= 10 {
		fmt.Println("[GenerateApiToken] Token limit exceeded")
		return "", errors.ErrTokenLimitExceeded
	}

	// Генерируем случайный токен
	token, err := randstring.GenerateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Хешируем токен для сохранения в БД
	tokenHash := hash.HashToken(token)

	// Создаём токен в БД
	_, err = s.apiTokenRepo.CreateApiToken(ctx, userID, subscription.ID, tokenName, tokenHash)
	if err != nil {
		return "", fmt.Errorf("failed to create API token: %w", err)
	}

	// Логируем создание токена
	s.auditLogRepo.CreateAuditLog(
		ctx,
		userID,
		"token_created",
		"api_token",
		0,
		nil,
		200,
		"",
		"",
		"",
	)

	// Возвращаем обычный токен (он будет отправлен клиенту только один раз)
	return token, nil
}

// ValidateApiToken проверяет валидность API токена
func (s *tokenService) ValidateApiToken(ctx context.Context, token string) (*sqlc.ApiToken, error) {
	// Хешируем токен
	tokenHash := hash.HashToken(token)

	// Получаем токен из БД
	apiToken, err := s.apiTokenRepo.GetApiTokenByHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get API token: %w", err)
	}

	// Проверяем, активен ли токен
	if !apiToken.IsActive.Bool {
		return nil, errors.ErrInvalidToken
	}

	// Проверяем, не истёк ли токен
	if apiToken.ExpiresAt.Valid && apiToken.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.ErrInvalidToken
	}

	// Проверяем, не был ли токен отозван
	if apiToken.RevokedAt.Valid {
		return nil, errors.ErrInvalidToken
	}

	// Обновляем время последнего использования
	s.apiTokenRepo.UpdateApiTokenLastUsed(ctx, apiToken.ID)

	return apiToken, nil
}

// RevokeApiToken отзывает API токен
func (s *tokenService) RevokeApiToken(ctx context.Context, tokenID int32) error {
	fmt.Printf("[RevokeApiToken] Revoking token ID: %d\n", tokenID)

	token, err := s.apiTokenRepo.RevokeApiToken(ctx, tokenID)
	if err != nil {
		fmt.Printf("[RevokeApiToken] Error revoking token: %v\n", err)
		return fmt.Errorf("failed to revoke API token: %w", err)
	}

	fmt.Printf("[RevokeApiToken] Token revoked successfully: ID=%d, UserID=%d\n", token.ID, token.UserID)

	return nil
}

// ListUserTokens получает все токены пользователя
func (s *tokenService) ListUserTokens(ctx context.Context, userID int32) ([]sqlc.ApiToken, error) {
	tokens, err := s.apiTokenRepo.ListUserApiTokens(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API tokens: %w", err)
	}

	return tokens, nil
}

// GetTokenByID получает токен по ID
func (s *tokenService) GetTokenByID(ctx context.Context, tokenID int32) (*sqlc.ApiToken, error) {
	token, err := s.apiTokenRepo.GetApiTokenByID(ctx, tokenID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}
