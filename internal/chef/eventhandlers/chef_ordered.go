
package eventhandlers

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/chef/domain"
	"github.com/ntv97/atriaseniorliving/internal/chef/infras/postgresql"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type chefOrderedEventHandler struct {
	pg         postgres.DBEngine
	waitstaffPub pkgPublisher.EventPublisher
}

var _ ChefOrderedEventHandler = (*chefOrderedEventHandler)(nil)

var ChefOrderedEventHandlerSet = wire.NewSet(NewChefOrderedEventHandler)

func NewChefOrderedEventHandler(
	pg postgres.DBEngine,
	waitstaffPub pkgPublisher.EventPublisher,
) ChefOrderedEventHandler {
	return &chefOrderedEventHandler{
		pg:         pg,
		waitstaffPub: waitstaffPub,
	}
}

func (h *chefOrderedEventHandler) Handle(ctx context.Context, e event.ChefOrdered) error {
	slog.Info("chefOrderedEventHandler-Handle", "ChefOrdered", e)

	order := domain.NewChefOrder(e)

	db := h.pg.GetDB()
	querier := postgresql.New(db)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "chefOrderedEventHandler.Handle")
	}

	qtx := querier.WithTx(tx)

	_, err = qtx.CreateOrder(ctx, postgresql.CreateOrderParams{
		ID:       order.ID,
		OrderID:  e.OrderID,
		ItemType: int32(order.ItemType),
		ItemName: order.ItemName,
		TimeUp:   order.TimeUp,
		Created:  order.Created,
		Updated: sql.NullTime{
			Time:  order.Updated,
			Valid: true,
		},
	})
	if err != nil {
		slog.Info("failed to call to repo", "error", err)

		return errors.Wrap(err, "chefOrderedEventHandler-querier.CreateOrder")
	}

	// todo: it might cause dual-write problem, but we accept it temporary
	for _, event := range order.DomainEvents() {
		eventBytes, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "json.Marshal[event]")
		}

		if err := h.waitstaffPub.Publish(ctx, eventBytes, "text/plain"); err != nil {
			return errors.Wrap(err, "waitstaffPub.Publish")
		}
	}

	return tx.Commit()
}

