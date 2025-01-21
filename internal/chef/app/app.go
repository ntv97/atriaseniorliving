
package app

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ntv97/atriaseniorliving/cmd/chef/config"
	"github.com/ntv97/atriaseniorliving/internal/chef/eventhandlers"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	pkgConsumer "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type App struct {
	Cfg *config.Config

	PG       postgres.DBEngine
	AMQPConn *amqp.Connection

	WaitstaffOrderPub pkgPublisher.EventPublisher
	Consumer        pkgConsumer.EventConsumer

	handler eventhandlers.ChefOrderedEventHandler
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	amqpConn *amqp.Connection,
	waitstaffOrderPub pkgPublisher.EventPublisher,
	consumer pkgConsumer.EventConsumer,
	handler eventhandlers.ChefOrderedEventHandler,
) *App {
	return &App{
		Cfg:      cfg,
		PG:       pg,
		AMQPConn: amqpConn,

		WaitstaffOrderPub: waitstaffOrderPub,
		Consumer:        consumer,

		handler: handler,
	}
}

func (c *App) Worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "chef-order-created":
			var payload event.ChefOrdered
			err := json.Unmarshal(delivery.Body, &payload)

			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = c.handler.Handle(ctx, payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					slog.Error("failed to delivery.Reject", err)
				}

				slog.Error("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					slog.Error("failed to acknowledge delivery", err)
				}
			}
		default:
			slog.Info("default")
		}
	}

	slog.Info("deliveries channel closed")
}

