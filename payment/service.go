package payment

import (
	"context"
	"santapan_payment_service/domain"
)

// PaymentService represents the payment service interface
type PostgresPaymentCommandRepository interface {
	Store(ctx context.Context, payment *domain.Payment) error
	Update(ctx context.Context, payment *domain.Payment) error
}

type PostgresPaymentQueryRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Payment, error)
	GetByRefID(ctx context.Context, refID string) (domain.Payment, error)
}

// PaymentHandler represents the HTTP handler for payments
type Service struct {
	postgresRepoQuery   PostgresPaymentQueryRepository
	postgresRepoCommand PostgresPaymentCommandRepository
}

// NewPaymentHandler initializes the payment endpoints
func NewService(postgresRepoQuery PostgresPaymentQueryRepository, postgresRepoCommand PostgresPaymentCommandRepository) *Service {
	return &Service{
		postgresRepoQuery:   postgresRepoQuery,
		postgresRepoCommand: postgresRepoCommand,
	}
}

// ProcessPayment processes a new payment
func (h *Service) Store(ctx context.Context, payment *domain.Payment) error {
	return h.postgresRepoCommand.Store(ctx, payment)
}

// UpdatePayment updates a payment
func (h *Service) Update(ctx context.Context, payment *domain.Payment) error {
	return h.postgresRepoCommand.Update(ctx, payment)
}

// GetPaymentByID retrieves a payment by its ID
func (h *Service) GetByID(ctx context.Context, id int64) (domain.Payment, error) {
	return h.postgresRepoQuery.GetByID(ctx, id)
}

// CallbackPayment processes a payment callback
func (h *Service) GetByRefID(ctx context.Context, refID string) (domain.Payment, error) {
	return h.postgresRepoQuery.GetByRefID(ctx, refID)
}
