package otp

import (
	"context"
	"santapan/domain"
	"time"
)

type PostgresRepositoryQueries interface {
	GetByID(ctx context.Context, id int64) (res domain.Otp, err error)
	GetByUserID(ctx context.Context, id int64) (res domain.Otp, err error)
}

type PostgresRepositoryCommand interface {
	Store(ctx context.Context, otp *domain.Otp) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, otp *domain.Otp) (err error)
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

func (s *Service) Store(ctx context.Context, otp *domain.Otp) (err error) {
	err = s.postgresRepoCommand.Store(ctx, otp)
	return

}

func (s *Service) GetByID(ctx context.Context, id int64) (res domain.Otp, err error) {
	res, err = s.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (s *Service) GetByUserID(ctx context.Context, id int64) (res domain.Otp, err error) {
	res, err = s.postgresRepoQuery.GetByUserID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *Service) Update(ctx context.Context, otp *domain.Otp) (err error) {
	otp.UpdatedAt = time.Now()
	return a.postgresRepoCommand.Update(ctx, otp)
}

func (s *Service) Delete(ctx context.Context, id int64) (err error) {
	existedToken, err := s.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return
	}
	return s.postgresRepoCommand.Delete(ctx, existedToken.ID)
}
