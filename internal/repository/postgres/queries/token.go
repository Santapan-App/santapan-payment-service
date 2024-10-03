package queries

import (
	"context"
	"database/sql"
	"tobby/domain"

	"github.com/sirupsen/logrus"
)

type PostgresTokenQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresTokenQueryRepository(conn *sql.DB) *PostgresTokenQueryRepository {
	return &PostgresTokenQueryRepository{conn}
}

func (m *PostgresTokenQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Token, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error("Error querying database: ", err)
		return nil, err
	}
	defer func() {
		if errRow := rows.Close(); errRow != nil {
			logrus.Error("Error closing rows: ", errRow)
		}
	}()

	result = make([]domain.Token, 0)
	for rows.Next() {
		var t domain.Token
		if err = rows.Scan(
			&t.ID,
			&t.RefreshToken,
			&t.UserID,
			&t.UpdatedAt,
			&t.CreatedAt,
		); err != nil {
			logrus.Error("Error scanning row: ", err)
			return nil, err
		}
		result = append(result, t)
	}

	if err = rows.Err(); err != nil {
		logrus.Error("Error iterating over rows: ", err)
		return nil, err
	}

	return result, nil
}

func (m *PostgresTokenQueryRepository) GetByID(ctx context.Context, id int64) (res domain.Token, err error) {
	query := `SELECT id, refresh_token, user_id, device_id, updated_at, created_at FROM token WHERE id = $1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Token{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *PostgresTokenQueryRepository) GetByUserID(ctx context.Context, id int64) (res domain.Token, err error) {
	query := `SELECT id, refresh_token, user_id, updated_at, created_at FROM token WHERE user_id = $1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Token{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
