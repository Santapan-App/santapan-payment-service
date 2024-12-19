package commands

import (
	"context"
	"database/sql"
	"fmt"
	"santapan_payment_service/domain"
)

// PostgresPaymentCommandRepository struct
type PostgresPaymentCommandRepository struct {
	Conn *sql.DB
}

// NewPostgresPaymentCommandRepository creates a new instance of PostgresPaymentCommandRepository
func NewPostgresPaymentCommandRepository(conn *sql.DB) *PostgresPaymentCommandRepository {
	return &PostgresPaymentCommandRepository{Conn: conn}
}

// Store stores the payment details
func (r *PostgresPaymentCommandRepository) Store(ctx context.Context, payment *domain.Payment) (err error) {
	query := `INSERT INTO payment (reference_id, session_id, user_id, amount, status, url, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
              RETURNING id`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, payment.ReferenceID, payment.SessionID, payment.UserID, payment.Amount, payment.Status, payment.Url, payment.CreatedAt, payment.UpdatedAt).Scan(&payment.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Update updates the payment status based on the payment ID
func (r *PostgresPaymentCommandRepository) Update(ctx context.Context, payment *domain.Payment) (err error) {
	query := `UPDATE payment SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, payment.Status, payment.ID)
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
