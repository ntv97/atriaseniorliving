
package infras

import (
	"context"

	"github.com/google/wire"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	"github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
)

var (
	CookEventPublisherSet = wire.NewSet(NewCookEventPublisher)
	ChefEventPublisherSet = wire.NewSet(NewChefEventPublisher)
	WaitstaffEventPublisherSet = wire.NewSet(NewWaitstaffEventPublisher)
)

type (
	cookEventPublisher struct {
		pub publisher.EventPublisher
	}
	chefEventPublisher struct {
		pub publisher.EventPublisher
	}
	waitstaffEventPublisher struct {
                pub publisher.EventPublisher
        }
)

func NewCookEventPublisher(pub publisher.EventPublisher) orders.CookEventPublisher {
	return &cookEventPublisher{
		pub: pub,
	}
}

func (p *cookEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

func (p *cookEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

func NewChefEventPublisher(pub publisher.EventPublisher) orders.ChefEventPublisher {
	return &chefEventPublisher{
		pub: pub,
	}
}

func (p *chefEventPublisher) Configure(opts ...publisher.Option) {
	p.pub.Configure(opts...)
}

func (p *chefEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return p.pub.Publish(ctx, body, contentType)
}

func NewWaitstaffEventPublisher(pub publisher.EventPublisher) orders.WaitstaffEventPublisher {
        return &waitstaffEventPublisher{
                pub: pub,
        }
}

func (p *waitstaffEventPublisher) Configure(opts ...publisher.Option) {
        p.pub.Configure(opts...)
}

func (p *waitstaffEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
        return p.pub.Publish(ctx, body, contentType)
}
