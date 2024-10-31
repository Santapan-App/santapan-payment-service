package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type PostgresDeviceQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresDeviceQueryRepository(conn *sql.DB) *PostgresDeviceQueryRepository {
	return &PostgresDeviceQueryRepository{conn}
}

// fetch is a helper method to execute a query and map the results to a slice of domain.Device
func (m *PostgresDeviceQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Device, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	for rows.Next() {
		t := domain.Device{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Brand,
			&t.UniqueID,
			&t.UserID,
			&t.IPAddress,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// GetByUserID retrieves a single device associated with a specific user ID
func (m *PostgresDeviceQueryRepository) GetByUserID(ctx context.Context, userID int64) (result domain.Device, err error) {
	query := `SELECT id, name, brand, unique_id, user_id, ip_address, created_at, updated_at, deleted_at
              FROM devices 
              WHERE user_id = $1
              LIMIT 1`

	devices, err := m.fetch(ctx, query, userID)
	if err != nil {
		return result, err
	}

	// Check if any devices were returned
	if len(devices) > 0 {
		return devices[0], nil // Return the first device found
	}

	return result, nil // No device found, return zero value
}
