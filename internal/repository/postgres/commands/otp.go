package commands

import (
	"context"
	"database/sql"
	"fmt"
	"santapan/domain"
)

type PostgresOtpCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresOtpCommandRepository(Conn *sql.DB) *PostgresOtpCommandRepository {
	return &PostgresOtpCommandRepository{Conn}
}

func (r *PostgresOtpCommandRepository) Store(ctx context.Context, otp *domain.Otp) (err error) {
	query := `INSERT otp SET code=?, device_id=?, user_id=?, created_at=?, updated_at=?`
	stmt, err := r.Conn.PrepareContext(ctx, query)

	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, otp.Code, otp.DeviceId, otp.UserId, otp.CreatedAt, otp.UpdatedAt)

	lastID, err := res.LastInsertId()

	if err != nil {
		return
	}

	otp.ID = lastID

	return
}

func (m *PostgresOtpCommandRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM otp WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (m *PostgresOtpCommandRepository) Update(ctx context.Context, otp *domain.Otp) (err error) {
	query := `UPDATE otp SET code=? , user_id=? , device_id=?, updated_at=? , created_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, otp.Code, otp.UserId, otp.DeviceId, otp.UpdatedAt, otp.CreatedAt)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
