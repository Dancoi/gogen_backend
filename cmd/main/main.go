package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Dancoi/gogen_backend/config"
	"github.com/Dancoi/gogen_backend/internal/handlers"
	"github.com/Dancoi/gogen_backend/internal/middleware"
	"github.com/Dancoi/gogen_backend/internal/repository"
	"github.com/Dancoi/gogen_backend/internal/service"
	"github.com/Dancoi/gogen_backend/internal/sqlc"
)

func main() {
	// Загружаем конфиг
	cfg := config.LoadConfig()
	fmt.Println("Config loaded:", cfg)

	// Подключаемся к БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверяем коннекшен
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Инициализируем repositories
	userRepo := repository.NewUserRepository(pool)
	subscriptionRepo := repository.NewSubscriptionRepository(pool)
	subscriptionPlanRepo := repository.NewSubscriptionPlanRepository(pool)
	apiTokenRepo := repository.NewApiTokenRepository(pool)
	userSessionRepo := repository.NewUserSessionRepository(pool)
	auditLogRepo := repository.NewAuditLogRepository(pool)

	// Инициализируем services
	authService := service.NewAuthService(
		userRepo,
		subscriptionRepo,
		subscriptionPlanRepo,
		userSessionRepo,
		auditLogRepo,
	)

	tokenService := service.NewTokenService(
		apiTokenRepo,
		subscriptionRepo,
		auditLogRepo,
	)

	// Инициализируем subscription plans
	if err := initSubscriptionPlans(ctx, pool); err != nil {
		log.Printf("Warning: Failed to initialize subscription plans: %v", err)
	}

	// Инициализируем handlers
	authHandler := handlers.NewAuthHandler(authService)
	tokenHandler := handlers.NewTokenHandler(tokenService)

	// Инициализируем Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes (без middleware)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// API routes (с middleware)
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(authService))
	{
		apiGroup.POST("/tokens", tokenHandler.GenerateToken)
		apiGroup.GET("/tokens", tokenHandler.ListTokens)
		apiGroup.DELETE("/tokens/:id", tokenHandler.RevokeToken)
		apiGroup.POST("/logout", authHandler.Logout)
	}

	fmt.Println("Starting server on port", cfg.ServerPort)

	// Запускаем сервер
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initSubscriptionPlans инициализирует планы подписки если их нет
func initSubscriptionPlans(ctx context.Context, pool *pgxpool.Pool) error {
	queries := sqlc.New(pool)

	// Проверяем есть ли trial план
	_, err := queries.GetSubscriptionPlanByType(ctx, "trial")
	if err == nil {
		// Планы уже существуют
		fmt.Println("Subscription plans already initialized")
		return nil
	}

	fmt.Println("Initializing subscription plans...")

	// Trial план
	trialParams := sqlc.CreateSubscriptionPlanParams{
		Name:              "Trial",
		PlanType:          "trial",
		MaxTokensPerMonth: 100,
		MaxApiCallsPerDay: 1000,
		Price:             pgtype.Numeric{Int: big.NewInt(0), Valid: true},
		TrialDurationDays: pgtype.Int4{Int32: 30, Valid: true},
		IsActive:          pgtype.Bool{Bool: true, Valid: true},
		Features:          []byte(`{"console_tool":true,"api_access":true,"support":"email"}`),
	}

	_, err = queries.CreateSubscriptionPlan(ctx, trialParams)
	if err != nil {
		return fmt.Errorf("failed to create trial plan: %w", err)
	}
	fmt.Println("✓ Trial plan created")

	// Premium план
	premiumParams := sqlc.CreateSubscriptionPlanParams{
		Name:              "Premium",
		PlanType:          "premium",
		MaxTokensPerMonth: 1000,
		MaxApiCallsPerDay: 10000,
		Price:             pgtype.Numeric{Int: big.NewInt(2999), Valid: true},
		TrialDurationDays: pgtype.Int4{Int32: 0, Valid: true},
		IsActive:          pgtype.Bool{Bool: true, Valid: true},
		Features:          []byte(`{"console_tool":true,"api_access":true,"support":"priority"}`),
	}

	_, err = queries.CreateSubscriptionPlan(ctx, premiumParams)
	if err != nil {
		return fmt.Errorf("failed to create premium plan: %w", err)
	}
	fmt.Println("✓ Premium plan created")

	// Commercial план
	commercialParams := sqlc.CreateSubscriptionPlanParams{
		Name:              "Commercial",
		PlanType:          "commercial",
		MaxTokensPerMonth: 10000,
		MaxApiCallsPerDay: 100000,
		Price:             pgtype.Numeric{Int: big.NewInt(29999), Valid: true},
		TrialDurationDays: pgtype.Int4{Int32: 0, Valid: true},
		IsActive:          pgtype.Bool{Bool: true, Valid: true},
		Features:          []byte(`{"console_tool":true,"api_access":true,"support":"24/7"}`),
	}

	_, err = queries.CreateSubscriptionPlan(ctx, commercialParams)
	if err != nil {
		return fmt.Errorf("failed to create commercial plan: %w", err)
	}
	fmt.Println("✓ Commercial plan created")

	fmt.Println("All subscription plans initialized successfully!")
	return nil
}
