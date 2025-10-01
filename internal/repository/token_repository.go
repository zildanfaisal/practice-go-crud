package repository

import (
	"context"
	"database/sql"
	"time"
)

type TokenRepository interface {
	AddToBlacklist(ctx context.Context, jti string, userID int, expiredAt time.Time) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
	CleanupExpired(ctx context.Context) error
}

type mysqlTokenRepo struct {
	db *sql.DB
}

func NewMySQLTokenRepo(db *sql.DB) TokenRepository {
	return &mysqlTokenRepo{db: db}
}

func (r *mysqlTokenRepo) AddToBlacklist(ctx context.Context, jti string, userID int, expiredAt time.Time) error {
	query := `INSERT INTO token_blacklist (jti, user_id, expired_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, jti, userID, expiredAt)
	return err
}

func (r *mysqlTokenRepo) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = ? AND expired_at > NOW())`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&exists)
	return exists, err
}

func (r *mysqlTokenRepo) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM token_blacklist WHERE expired_at <= NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
