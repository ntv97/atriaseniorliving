
package eventhandlers

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/cook/domain"
	"github.com/ntv97/atriaseniorliving/internal/cook/infras/postgresql"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
	pkgPublisher "github.com/ntv97/atriaseniorliving/pkg/rabbitmq/publisher"
	"golang.org/x/exp/slog"
)

type cookOrderedEventHandler struct {
	pg         postgres.DBEngine
	waitstaffPub pkgPublisher.EventPublisher
}

var _ CookOrderedEventHandler = (*cookOrderedEventHandler)(nil)

var CookOrderedEventHandlerSet = wire.NewSet(NewCookOrderedEventHandler)

func NewCookOrderedEventHandler(
	pg postgres.DBEngine,
	waitstaffPub pkgPublisher.EventPublisher,
) CookOrderedEventHandler {
	return &cookOrderedEventHandler{
		pg:         pg,
		waitstaffPub: waitstaffPub,
	}
}

func (h *cookOrderedEventHandler) Handle(ctx context.Context, e event.CookOrdered) error {
	slog.Info("cookOrderedEventHandler-Handle", "CookOrdered", e)

	order := domain.NewCookOrder(e)

	db := h.pg.GetDB()
	querier := postgresql.New(db)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "cookOrderedEventHandler.Handle")
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

		return errors.Wrap(err, "cookOrderedEventHandler-querier.CreateOrder")
	}

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

