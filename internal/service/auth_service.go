package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Dancoi/gogen_backend/internal/repository"
	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/Dancoi/gogen_backend/pkg/utils/errors"
	"github.com/Dancoi/gogen_backend/pkg/utils/jwt"
	passwordutil "github.com/Dancoi/gogen_backend/pkg/utils/password"
	"github.com/Dancoi/gogen_backend/pkg/utils/randstring"
)

// AuthService интерфейс для аутентификации
type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*sqlc.User, error)
	Login(ctx context.Context, email, password string) (*sqlc.User, string, error)
	ValidateToken(tokenString string) (*jwt.Claims, error)
}

type authService struct {
	userRepo         repository.UserRepository
	subscriptionRepo repository.SubscriptionRepository
	subscriptionPlan repository.SubscriptionPlanRepository
	userSessionRepo  repository.UserSessionRepository
	auditLogRepo     repository.AuditLogRepository
}

// NewAuthService создаёт новый AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	subscriptionRepo repository.SubscriptionRepository,
	subscriptionPlan repository.SubscriptionPlanRepository,
	userSessionRepo repository.UserSessionRepository,
	auditLogRepo repository.AuditLogRepository,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		subscriptionPlan: subscriptionPlan,
		userSessionRepo:  userSessionRepo,
		auditLogRepo:     auditLogRepo,
	}
}

// Register регистрирует нового пользователя
func (s *authService) Register(ctx context.Context, email, username, password string) (*sqlc.User, error) {
	fmt.Printf("[Register] Starting registration for email: %s\n", email)

	// Проверяем, существует ли пользователь с таким email
	_, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil {
		fmt.Printf("[Register] Email already exists: %s\n", email)
		return nil, errors.ErrEmailAlreadyExists
	}
	if err.Error() != "no rows in result set" {
		fmt.Printf("[Register] Error checking email: %v\n", err)
		return nil, err
	}
	fmt.Println("[Register] Email is available")

	// Проверяем, существует ли пользователь с таким username
	_, err = s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		fmt.Printf("[Register] Username already exists: %s\n", username)
		return nil, errors.ErrUsernameAlreadyExists
	}
	if err.Error() != "no rows in result set" {
		fmt.Printf("[Register] Error checking username: %v\n", err)
		return nil, err
	}
	fmt.Println("[Register] Username is available")

	// Хешируем пароль
	passwordHash, err := passwordutil.HashPassword(password)
	if err != nil {
		fmt.Printf("[Register] Error hashing password: %v\n", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	fmt.Println("[Register] Password hashed")

	// Создаём пользователя
	user, err := s.userRepo.CreateUser(ctx, email, username, passwordHash)
	if err != nil {
		fmt.Printf("[Register] Error creating user: %v\n", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Printf("[Register] User created with ID: %s\n", user.ID)

	// Создаём trial подписку
	fmt.Println("[Register] Fetching trial plan...")
	trialPlan, err := s.subscriptionPlan.GetSubscriptionPlanByType(ctx, "trial")
	if err != nil {
		fmt.Printf("[Register] Error fetching trial plan: %v\n", err)
		// Если плана нет, продолжаем без подписки (лучше логировать)
		return user, nil
	}
	fmt.Printf("[Register] Trial plan found: ID=%s, Duration=%d days\n", trialPlan.ID, trialPlan.TrialDurationDays.Int32)

	trialDays := trialPlan.TrialDurationDays
	if trialPlan.TrialDurationDays.Int32 > 0 {
		trialDays = trialPlan.TrialDurationDays
	}

	expiresAt := time.Now().AddDate(0, 0, int(trialDays.Int32)).Unix()
	fmt.Printf("[Register] Creating subscription with expiration at: %v\n", time.Unix(expiresAt, 0))

	_, err = s.subscriptionRepo.CreateSubscription(
		ctx,
		user.ID,
		trialPlan.ID,
		"trial",
		expiresAt,
		trialPlan.MaxTokensPerMonth,
	)
	if err != nil {
		fmt.Printf("[Register] Error creating subscription: %v\n", err)
		return nil, fmt.Errorf("failed to create trial subscription: %w", err)
	}
	fmt.Println("[Register] Subscription created")

	// Логируем регистрацию
	s.auditLogRepo.CreateAuditLog(
		ctx,
		user.ID,
		"user_created",
		"user",
		user.ID,
		nil,
		200,
		"",
		"",
		"",
	)
	fmt.Println("[Register] Registration completed successfully")

	return user, nil
}

// Login авторизует пользователя и возвращает JWT токен
func (s *authService) Login(ctx context.Context, email, password string) (*sqlc.User, string, error) {
	// Получаем пользователя по email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			s.auditLogRepo.CreateAuditLog(ctx, 0, "login_failed", "user", 0, nil, 401, "user not found", "", "")
			return nil, "", errors.ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Проверяем пароль
	if !passwordutil.VerifyPassword(user.PasswordHash, password) {
		s.auditLogRepo.CreateAuditLog(ctx, user.ID, "login_failed", "user", user.ID, nil, 401, "invalid password", "", "")
		return nil, "", errors.ErrInvalidCredentials
	}

	// Генерируем session token
	sessionToken, err := randstring.GenerateRandomString(32)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate session token: %w", err)
	}

	// Создаём сессию на 24 часа
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	session, err := s.userSessionRepo.CreateUserSession(ctx, user.ID, sessionToken, "", "", expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create session: %w", err)
	}

	// Генерируем JWT токен
	jwtToken, err := jwt.GenerateToken(user.ID, user.Email, user.Username, session.ID, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Логируем успешный вход
	s.auditLogRepo.CreateAuditLog(ctx, user.ID, "login", "user", user.ID, nil, 200, "", "", "")

	return user, jwtToken, nil
}

// ValidateToken проверяет JWT токен
func (s *authService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return jwt.VerifyToken(tokenString)
}
