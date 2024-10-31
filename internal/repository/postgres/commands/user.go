package commands

import (
	"context"
	"database/sql"
	"fmt"
	"santapan/domain"
	"time"
)

type PostgresUserCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresUserCommandRepository(conn *sql.DB) *PostgresUserCommandRepository {
	return &PostgresUserCommandRepository{Conn: conn}
}

func (r *PostgresUserCommandRepository) Store(ctx context.Context, user *domain.User) (err error) {
	query := `INSERT INTO users (full_name, email, password, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user.FullName, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (r *PostgresUserCommandRepository) UpdatePhoneVerifiedAt(ctx context.Context, id int64, emailVerifiedAt time.Time) error {
	query := `UPDATE users SET email_verified_at=$1 WHERE id=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, emailVerifiedAt, id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affect != 1 {
		return fmt.Errorf("unexpected number of affected rows: %d", affect)
	}

	return nil
}
