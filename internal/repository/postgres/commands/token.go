package commands

import (
	"context"
	"database/sql"
	"fmt"
	"tobby/domain"
)

type PostgresTokenCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresTokenCommandRepository(Conn *sql.DB) *PostgresTokenCommandRepository {
	return &PostgresTokenCommandRepository{Conn}
}

func (r *PostgresTokenCommandRepository) Store(ctx context.Context, token *domain.Token) (err error) {
	query := `INSERT INTO token (refresh_token, user_id, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, token.RefreshToken, token.UserID).Scan(&token.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (m *PostgresTokenCommandRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM token WHERE id = $1"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("unexpected number of affected rows: %d", rowsAffected)
	}

	return nil
}

func (m *PostgresTokenCommandRepository) Update(ctx context.Context, token *domain.Token) (err error) {
	query := `UPDATE token SET refresh_token = $1, user_id = $2 WHERE id = $3`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare update query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, token.RefreshToken, token.UserID, token.ID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("unexpected number of affected rows: %d", rowsAffected)
	}

	return nil
}
