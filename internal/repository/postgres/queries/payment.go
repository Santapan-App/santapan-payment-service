package queries

import (
	"context"
	"database/sql"
	"fmt"
	"santapan_payment_service/domain"

	"github.com/sirupsen/logrus"
)

// PostgresPaymentQueryRepository struct
type PostgresPaymentQueryRepository struct {
	conn *sql.DB
}

// NewPostgresPaymentQueryRepository creates a new instance of PostgresPaymentQueryRepository
func NewPostgresPaymentQueryRepository(conn *sql.DB) *PostgresPaymentQueryRepository {
	return &PostgresPaymentQueryRepository{conn: conn}
}

// Helper function to execute a generic query
func (m *PostgresPaymentQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Payment, err error) {
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Error("Failed to execute query")
		return nil, err
	}

	defer func() {
		if errRow := rows.Close(); errRow != nil {
			logrus.Error("Failed to close rows:", errRow)
		}
	}()

	result = make([]domain.Payment, 0)
	for rows.Next() {
		payment := domain.Payment{}
		err = rows.Scan(
			&payment.ID,
			&payment.ReferenceID,
			&payment.SessionID,
			&payment.UserID,
			&payment.Amount,
			&payment.Status,
			&payment.Url,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			logrus.Error("Failed to scan row:", err)
			return nil, err
		}
		result = append(result, payment)
	}

	return result, nil
}

// GetByUserID retrieves payments associated with a specific user ID
func (m *PostgresPaymentQueryRepository) GetByUserID(ctx context.Context, userID int64) (res []domain.Payment, err error) {
	query := `SELECT id, reference_id, session_id, user_id, amount, status, url, created_at, updated_at FROM payment WHERE user_id = $1`
	return m.fetch(ctx, query, userID)
}

// GetByID retrieves a payment by its ID
func (m *PostgresPaymentQueryRepository) GetByID(ctx context.Context, id int64) (res domain.Payment, err error) {
	query := `SELECT id, reference_id, session_id, user_id, amount, status, url, created_at, updated_at FROM payment WHERE id = $1`
	result, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Payment{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return domain.Payment{}, sql.ErrNoRows
}

// GetByRefID
func (m *PostgresPaymentQueryRepository) GetByRefID(ctx context.Context, refID string) (res domain.Payment, err error) {
	query := `SELECT id, reference_id, session_id, user_id, amount, status, url, created_at, updated_at FROM payment WHERE reference_id = $1`
	result, err := m.fetch(ctx, query, refID)
	if err != nil {
		return domain.Payment{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return domain.Payment{}, sql.ErrNoRows
}

// Validate checks if the payment exists for a given ID and user ID
func (m *PostgresPaymentQueryRepository) Validate(ctx context.Context, paymentID int64, userID int64) (err error) {
	query := `SELECT id FROM payment WHERE id = $1 AND user_id = $2`
	result, err := m.fetch(ctx, query, paymentID, userID)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return fmt.Errorf("payment not found for the user")
	}

	return nil
}
