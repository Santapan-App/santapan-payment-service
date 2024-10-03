package queries

import (
	"context"
	"database/sql"
	"tobby/domain"

	"github.com/sirupsen/logrus"
)

type PostgresUserQueryRepository struct {
	conn *sql.DB
}

func NewPostgresUserQueryRepository(conn *sql.DB) *PostgresUserQueryRepository {
	return &PostgresUserQueryRepository{conn: conn}
}

func (m *PostgresUserQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
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

	result = make([]domain.User, 0)
	for rows.Next() {
		user := domain.User{}
		err = rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.EmailVerifiedAt,
		)

		if err != nil {
			logrus.Error("Failed to scan row:", err)
			return nil, err
		}
		result = append(result, user)
	}

	return result, nil
}

// GetByEmail retrieves a user by their phone number and country code
func (m *PostgresUserQueryRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT id, full_name, email, password, created_at, updated_at, deleted_at, email_verified_at FROM users WHERE email = $1`
	list, err := m.fetch(ctx, query, email)

	if err != nil {
		return domain.User{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
