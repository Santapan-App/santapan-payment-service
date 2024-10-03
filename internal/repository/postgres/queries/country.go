package queries

import (
	"context"
	"database/sql"
	"fmt"
	"tobby/domain"
	"tobby/internal/repository"

	"github.com/sirupsen/logrus"
)

type PostgresCountryQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresCountryQueryRepository(conn *sql.DB) *PostgresCountryQueryRepository {
	return &PostgresCountryQueryRepository{Conn: conn}
}

func (r *PostgresCountryQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (res []domain.Country, err error) {
	rows, err := r.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		country := domain.Country{}
		err = rows.Scan(&country.ID, &country.Code, &country.Name, &country.Phone, &country.CreatedAt, &country.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		res = append(res, country)
	}

	return
}

// GetByCode fetches a country by its code
func (r *PostgresCountryQueryRepository) GetByCode(ctx context.Context, code string) (result domain.Country, err error) {
	// Use fetch Function
	query := `SELECT id, code, name, phone, created_at, updated_at FROM countries WHERE code = $1`
	country, err := r.fetch(ctx, query, code)
	if err != nil {
		logrus.Error(err)
		return result, err
	}
	logrus.Info(country, code)
	// Check if any devices were returned
	if len(country) > 0 {
		return country[0], nil // Return the first device found
	}

	return result, nil // No device found, return zero value
}

// Fetch retrieves all countries
func (r *PostgresCountryQueryRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Country, nextCursor string, err error) {
	query := `SELECT id, code, name, phone, created_at, updated_at FROM countries WHERE id > $1 ORDER BY id LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)

	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = r.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor, err = repository.EncodeCursor(res[len(res)-1].ID)

		if err != nil {
			return nil, "", err
		}
	}

	return
}
