package repository

import (
	"context"

	"github.com/Dancoi/gogen_backend/internal/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type auditLogRepository struct {
	queries *sqlc.Queries
}

// NewAuditLogRepository создаёт новый AuditLogRepository
func NewAuditLogRepository(pool *pgxpool.Pool) AuditLogRepository {
	return &auditLogRepository{
		queries: sqlc.New(pool),
	}
}

// CreateAuditLog создаёт новый лог аудита
func (r *auditLogRepository) CreateAuditLog(ctx context.Context, userID int32, action, resourceType string, resourceID int32, changes []byte, statusCode int32, errorMessage, ipAddress, userAgent string) (*sqlc.AuditLog, error) {
	log, err := r.queries.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
		UserID:       pgtype.Int4{Int32: userID, Valid: true},
		Action:       action,
		ResourceType: pgtype.Text{String: resourceType, Valid: resourceType != ""},
		ResourceID:   pgtype.Int4{Int32: resourceID, Valid: resourceID != 0},
		Changes:      changes,
		StatusCode:   pgtype.Int4{Int32: statusCode, Valid: statusCode != 0},
		ErrorMessage: pgtype.Text{String: errorMessage, Valid: errorMessage != ""},
		IpAddress:    pgtype.Text{String: ipAddress, Valid: ipAddress != ""},
		UserAgent:    pgtype.Text{String: userAgent, Valid: userAgent != ""},
	})
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetUserAuditLogs получает логи аудита пользователя
func (r *auditLogRepository) GetUserAuditLogs(ctx context.Context, userID int32, limit int32, offset int32) ([]sqlc.AuditLog, error) {
	return r.queries.GetUserAuditLogs(ctx, sqlc.GetUserAuditLogsParams{
		UserID: pgtype.Int4{Int32: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})
}
