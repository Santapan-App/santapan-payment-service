package country

import (
	"context"
	"tobby/domain"
)

type PostgresRepositoryQueries interface {
	GetByCode(ctx context.Context, code string) (res domain.Country, err error)
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Country, nextCursor string, err error)
}

//go:generate mockery --name ArticleRepository
type Service struct {
	postgresRepoQuery PostgresRepositoryQueries
}

// NewService will create a new article service object
func NewService(pq PostgresRepositoryQueries) *Service {
	return &Service{
		postgresRepoQuery: pq,
	}
}

func (s *Service) GetByCode(ctx context.Context, code string) (res domain.Country, err error) {
	res, err = s.postgresRepoQuery.GetByCode(ctx, code)
	if err != nil {
		return domain.Country{}, err
	}
	return
}

func (s *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Country, nextCursor string, err error) {
	return s.postgresRepoQuery.Fetch(ctx, cursor, num)
}
