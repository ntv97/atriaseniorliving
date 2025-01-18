
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ntv97/atriaseniorliving/cmd/waitstaff/config"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/app/router"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/events/handlers"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/infras"
	infrasGRPC "github.com/ntv97/atriaseniorliving/internal/waitstaff/infras/grpc"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/infras/repo"
	ordersUC "github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	"github.com/ntv97/atriaseniorliving/pkg/rabbitmq"
	pkgConsumer "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		rabbitMQFunc,
		pkgPublisher.EventPublisherSet,
		pkgConsumer.EventConsumerSet,

		infras.BaristaEventPublisherSet,
		infras.KitchenEventPublisherSet,
		infras.WaitstaffEventPublisherSet,
		infrasGRPC.ItemsGRPCClientSet,
		router.WaitstaffGRPCServerSet,
		repo.RepositorySet,
		ordersUC.UseCaseSet,
		handlers.BaristaOrderUpdatedEventHandlerSet,
		handlers.KitchenOrderUpdatedEventHandlerSet,
		handlers.WaitstaffOrderUpdatedEventHandlerSet,
	))
}

func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}

func rabbitMQFunc(url rabbitmq.RabbitMQConnStr) (*amqp.Connection, func(), error) {
	conn, err := rabbitmq.NewRabbitMQConn(url)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}

