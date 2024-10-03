package commands

import (
	"context"
	"database/sql"
	"fmt"
	"tobby/domain"
)

type PostgresDeviceCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresDeviceCommandRepository(Conn *sql.DB) *PostgresDeviceCommandRepository {
	return &PostgresDeviceCommandRepository{Conn}
}

// Store inserts a new device into the database
func (r *PostgresDeviceCommandRepository) Store(ctx context.Context, device *domain.Device) (err error) {
	query := `INSERT INTO devices (name, brand, unique_id, user_id, ip_address, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	// Prepare the statement
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	err = stmt.QueryRowContext(ctx, device.Name, device.Brand, device.UniqueID, device.UserID, device.IPAddress, device.CreatedAt, device.UpdatedAt).Scan(&device.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Update modifies an existing device in the database
func (r *PostgresDeviceCommandRepository) Update(ctx context.Context, device *domain.Device) (err error) {
	query := `UPDATE devices 
              SET name = $1, brand = $2, unique_id = $3, ip_address = $4, updated_at = $5 
              WHERE id = $6`

	// Prepare the statement
	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.ExecContext(ctx, device.Name, device.Brand, device.UniqueID, device.IPAddress, device.UpdatedAt, device.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
