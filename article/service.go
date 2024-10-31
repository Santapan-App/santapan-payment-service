package article

import (
	"context"
	"santapan/domain"
)

type PostgresRepositoryQueries interface {
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error)
}

type PostgresRepositoryCommand interface {
}

//go:generate mockery --name ArticleRepository
type Service struct {
	postgresRepoQuery   PostgresRepositoryQueries
	postgresRepoCommand PostgresRepositoryCommand
}

// NewService will create a new article service object
func NewService(pq PostgresRepositoryQueries, pc PostgresRepositoryCommand) *Service {
	return &Service{
		postgresRepoQuery:   pq,
		postgresRepoCommand: pc,
	}
}

func (a *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	res, nextCursor, err = a.postgresRepoQuery.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

// GetByID

func (a *Service) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	res, err = a.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return
}
