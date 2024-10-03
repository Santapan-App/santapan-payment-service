package device

import (
	"context"
	"time"
	"tobby/domain"
)

type PostgresRepositoryQueries interface {
	GetByUserID(ctx context.Context, id int64) (res domain.Device, err error)
}

type PostgresRepositoryCommand interface {
	Store(ctx context.Context, token *domain.Device) (err error)
	Update(ctx context.Context, token *domain.Device) (err error)
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

func (s *Service) Store(ctx context.Context, rt *domain.Device) (err error) {
	err = s.postgresRepoCommand.Store(ctx, rt)
	return
}

func (s *Service) GetByUserID(ctx context.Context, id int64) (res domain.Device, err error) {
	res, err = s.postgresRepoQuery.GetByUserID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (a *Service) Update(ctx context.Context, device *domain.Device) (err error) {
	device.UpdatedAt = time.Now()
	return a.postgresRepoCommand.Update(ctx, device)
}
