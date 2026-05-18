package repository

import (
	"context"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository создаёт новый UserRepository
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{
		queries: sqlc.New(pool),
	}
}

// CreateUser создаёт нового пользователя
func (r *userRepository) CreateUser(ctx context.Context, email, username, passwordHash string) (*sqlc.User, error) {
	usr, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		IsActive:      pgtype.Bool{Bool: true, Valid: true},
		EmailVerified: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// GetUserByEmail получает пользователя по email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	usr, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// GetUserByID получает пользователя по ID
func (r *userRepository) GetUserByID(ctx context.Context, id int32) (*sqlc.User, error) {
	usr, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// GetUserByUsername получает пользователя по username
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error) {
	usr, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// UpdateUserPasswordHash обновляет хеш пароля
func (r *userRepository) UpdateUserPasswordHash(ctx context.Context, userID int32, passwordHash string) (*sqlc.User, error) {
	usr, err := r.queries.UpdateUserPasswordHash(ctx, sqlc.UpdateUserPasswordHashParams{
		ID:           userID,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// UpdateUserIsActive обновляет статус активности пользователя
func (r *userRepository) UpdateUserIsActive(ctx context.Context, userID int32, isActive bool) (*sqlc.User, error) {
	usr, err := r.queries.UpdateUserIsActive(ctx, sqlc.UpdateUserIsActiveParams{
		ID:       userID,
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// UpdateUserEmailVerified обновляет статус верификации email
func (r *userRepository) UpdateUserEmailVerified(ctx context.Context, userID int32, emailVerified bool) (*sqlc.User, error) {
	usr, err := r.queries.UpdateUserEmailVerified(ctx, sqlc.UpdateUserEmailVerifiedParams{
		ID:            userID,
		EmailVerified: pgtype.Bool{Bool: emailVerified, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// DeleteUser удаляет пользователя
func (r *userRepository) DeleteUser(ctx context.Context, userID int32) error {
	return r.queries.DeleteUser(ctx, userID)
}
