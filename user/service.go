package user

import (
	"context"
	"time"
	"tobby/domain"
)

// ArticleRepository represent the article's repository contract
//
//go:generate mockery --name ArticleRepository
type PostgresRepositoryQueries interface {
	GetByEmail(ctx context.Context, email string) (res domain.User, err error)
}

type PostgresRepositoryCommand interface {
	Store(ctx context.Context, user *domain.User) (err error)
	UpdatePhoneVerifiedAt(ctx context.Context, id int64, time time.Time) (err error)
}

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

func (s *Service) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	return s.postgresRepoQuery.GetByEmail(ctx, email)
}

func (s *Service) UpdatePhoneVerifiedAt(ctx context.Context, id int64, time time.Time) (err error) {
	return s.postgresRepoCommand.UpdatePhoneVerifiedAt(ctx, id, time)
}

func (s *Service) Store(ctx context.Context, user *domain.User) (err error) {
	return s.postgresRepoCommand.Store(ctx, user)
}
