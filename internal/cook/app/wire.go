
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ntv97/atriaseniorliving/cmd/cook/config"
	"github.com/ntv97/atriaseniorliving/internal/cook/eventhandlers"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	"github.com/ntv97/atriaseniorliving/pkg/rabbitmq"
	pkgConsumer "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	rabbitMQConnStr rabbitmq.RabbitMQConnStr,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		rabbitMQFunc,
		pkgPublisher.EventPublisherSet,
		pkgConsumer.EventConsumerSet,
		eventhandlers.CookOrderedEventHandlerSet,
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

