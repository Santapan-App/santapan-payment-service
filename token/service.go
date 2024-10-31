package token

import (
	"context"
	"santapan/domain"
	"time"
)

type PostgresRepositoryQueries interface {
	GetByID(ctx context.Context, id int64) (res domain.Token, err error)
	GetByUserID(ctx context.Context, id int64) (res domain.Token, err error)
}

type PostgresRepositoryCommand interface {
	Store(ctx context.Context, token *domain.Token) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, token *domain.Token) (err error)
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

func (s *Service) Store(ctx context.Context, rt *domain.Token) (err error) {
	err = s.postgresRepoCommand.Store(ctx, rt)
	return

}

func (s *Service) GetByID(ctx context.Context, id int64) (res domain.Token, err error) {
	res, err = s.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (s *Service) GetByUserID(ctx context.Context, id int64) (res domain.Token, err error) {
	res, err = s.postgresRepoQuery.GetByUserID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *Service) Update(ctx context.Context, token *domain.Token) (err error) {
	token.UpdatedAt = time.Now()
	return a.postgresRepoCommand.Update(ctx, token)
}

func (s *Service) Delete(ctx context.Context, id int64) (err error) {
	existedToken, err := s.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return
	}
	return s.postgresRepoCommand.Delete(ctx, existedToken.ID)
}
