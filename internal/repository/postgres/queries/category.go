package queries

import (
	"context"
	"database/sql"
	"santapan/domain"
	"santapan/internal/repository"

	"github.com/sirupsen/logrus"
)

type CategoryRepository struct {
	Conn *sql.DB
}

// NewCategoryRepository creates an instance of CategoryRepository.
func NewCategoryRepository(conn *sql.DB) *CategoryRepository {
	return &CategoryRepository{conn}
}

func (m *CategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Category, err error) {
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

	result = make([]domain.Category, 0)
	for rows.Next() {
		t := domain.Category{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.ImageURL,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// Fetch retrieves categories with ID-based pagination.
func (m *CategoryRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Category, nextCursor string, err error) {
	query := `SELECT id, title, image_url, created_at, updated_at
			  FROM category WHERE id > $1 ORDER BY id LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	// Set the nextCursor if the result count reaches the limit
	if len(res) == int(num) {
		nextCursor, err = repository.EncodeCursor(res[len(res)-1].ID)
		if err != nil {
			logrus.Error("Failed to encode cursor: ", err)
			return res, "", err
		}
	}

	return res, nextCursor, nil
}

// GetByID retrieves a category by its ID.
func (m *CategoryRepository) GetByID(ctx context.Context, id int64) (res domain.Category, err error) {
	query := `SELECT id, title, image_url, created_at, updated_at
			  FROM category WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Category{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

// GetByTitle retrieves a category by its title.
func (m *CategoryRepository) GetByTitle(ctx context.Context, title string) (res domain.Category, err error) {
	query := `SELECT id, title, image_url, created_at, updated_at
			  FROM category WHERE title = $1`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
