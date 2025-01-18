package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	"github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
)

type (
	OrderRepo interface {
		GetAll(context.Context) ([]*domain.Order, error)
		GetByID(context.Context, uuid.UUID) (*domain.Order, error)
		Create(context.Context, *domain.Order) error
		Update(context.Context, *domain.Order) (*domain.Order, error)
	}

	CookEventPublisher interface {
		Configure(...publisher.Option)
		Publish(context.Context, []byte, string) error
	}

	ChefEventPublisher interface {
		Configure(...publisher.Option)
		Publish(context.Context, []byte, string) error
	}

	WaitstaffEventPublisher interface {
                Configure(...publisher.Option)
                Publish(context.Context, []byte, string) error
        }

	UseCase interface {
		GetListOrderFulfillment(context.Context) ([]*domain.Order, error)
		PlaceOrder(context.Context, *domain.PlaceOrderModel) error
	}
)
