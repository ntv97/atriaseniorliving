package app

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ntv97/atriaseniorliving/cmd/waitstaff/config"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/events"
	ordersUC "github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/event"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	pkgConsumer "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
	"github.com/ntv97/atriaseniorliving/proto/gen"
	"golang.org/x/exp/slog"
)

type App struct {
	Cfg       *config.Config
	PG        postgres.DBEngine
	AMQPConn  *amqp.Connection
	Publisher pkgPublisher.EventPublisher
	Consumer  pkgConsumer.EventConsumer

	CookOrderPub ordersUC.CookEventPublisher
	ChefOrderPub ordersUC.ChefEventPublisher
	WaitstaffOrderPub ordersUC.WaitstaffEventPublisher

	ItemsDomainSvc  domain.ItemsDomainService
	UC                ordersUC.UseCase
	WaitstaffGRPCServer gen.WaitstaffServiceServer

	cookHandler events.CookOrderUpdatedEventHandler
	chefHandler events.ChefOrderUpdatedEventHandler
	waitstaffHandler events.WaitstaffOrderUpdatedEventHandler
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	amqpConn *amqp.Connection,
	publisher pkgPublisher.EventPublisher,
	consumer pkgConsumer.EventConsumer,

	cookOrderPub ordersUC.CookEventPublisher,
	chefOrderPub ordersUC.ChefEventPublisher,
	waitstaffOrderPub ordersUC.WaitstaffEventPublisher,
	itemsDomainSvc domain.ItemsDomainService,
	uc ordersUC.UseCase,
	waitstaffGRPCServer gen.WaitstaffServiceServer,

	cookHandler events.CookOrderUpdatedEventHandler,
	chefHandler events.ChefOrderUpdatedEventHandler,
	waitstaffHandler events.WaitstaffOrderUpdatedEventHandler,
) *App {
	return &App{
		Cfg: cfg,

		PG:        pg,
		AMQPConn:  amqpConn,
		Publisher: publisher,
		Consumer:  consumer,

		CookOrderPub: cookOrderPub,
		ChefOrderPub: chefOrderPub,
		WaitstaffOrderPub: waitstaffOrderPub,

		ItemsDomainSvc:  itemsDomainSvc,
		UC:                uc,
		WaitstaffGRPCServer: waitstaffGRPCServer,

		cookHandler: cookHandler,
		chefHandler: chefHandler,
		waitstaffHandler: waitstaffHandler,
	}
}

func (a *App) Worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		slog.Info("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		slog.Info("received", "delivery_type", delivery.Type)

		
		switch delivery.Type {
		case "cook-order-updated":
			var payload shared.CookOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = a.cookHandler.Handle(ctx, &payload)

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
		case "chef-order-updated":
			var payload shared.ChefOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				slog.Error("failed to Unmarshal message", err)
			}

			err = a.chefHandler.Handle(ctx, &payload)

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
		case "waitstaff-order-updated":
                        var payload shared.WaitstaffOrderUpdated

                        err := json.Unmarshal(delivery.Body, &payload)
                        if err != nil {
                                slog.Error("failed to Unmarshal message", err)
                        }

                        err = a.waitstaffHandler.Handle(ctx, &payload)

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
	slog.Info("Deliveries channel closed")
}
