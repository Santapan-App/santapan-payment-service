package queries

import (
	"context"
	"database/sql"
	"santapan/domain"
	"santapan/internal/repository"

	"github.com/sirupsen/logrus"
)

type ArticleRepository struct {
	Conn *sql.DB
}

// NewArticleRepository creates an instance of ArticleRepository.
func NewArticleRepository(conn *sql.DB) *ArticleRepository {
	return &ArticleRepository{conn}
}

func (m *ArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
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

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
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

// Fetch retrieves articles with ID-based pagination.
func (m *ArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id, title, content, image_url, created_at, updated_at
			  FROM article WHERE id > $1 ORDER BY id LIMIT $2`

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

// GetByID retrieves an article by its ID.
func (m *ArticleRepository) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	query := `SELECT id, title, content, image_url, created_at, updated_at
			  FROM article WHERE id = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

// GetByTitle retrieves an article by its title.
func (m *ArticleRepository) GetByTitle(ctx context.Context, title string) (res domain.Article, err error) {
	query := `SELECT id, title, content, image_url, created_at, updated_at
			  FROM article WHERE title = $1`

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
