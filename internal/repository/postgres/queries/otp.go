package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type PostgresOtpQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresOtpQueryRepository(conn *sql.DB) *PostgresOtpQueryRepository {
	return &PostgresOtpQueryRepository{conn}
}

func (m *PostgresOtpQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Otp, err error) {
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

	result = make([]domain.Otp, 0)
	for rows.Next() {
		t := domain.Otp{}
		err = rows.Scan(
			&t.ID,
			&t.Code,
			&t.Retry,
			&t.UserId,
			&t.DeviceId,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *PostgresOtpQueryRepository) GetByID(ctx context.Context, id int64) (res domain.Otp, err error) {
	query := `SELECT * FROM otp WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Otp{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *PostgresOtpQueryRepository) GetByUserID(ctx context.Context, id int64) (res domain.Otp, err error) {
	query := `SELECT * FROM otp WHERE user_id = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Otp{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
