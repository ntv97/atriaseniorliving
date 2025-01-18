
package handlers

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/events"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
)

type cookOrderUpdatedEventHandler struct {
	orderRepo orders.OrderRepo
}

var _ events.CookOrderUpdatedEventHandler = (*cookOrderUpdatedEventHandler)(nil)

var CookOrderUpdatedEventHandlerSet = wire.NewSet(NewCookOrderUpdatedEventHandler)

func NewCookOrderUpdatedEventHandler(orderRepo orders.OrderRepo) events.CookOrderUpdatedEventHandler {
	return &cookOrderUpdatedEventHandler{
		orderRepo: orderRepo,
	}
}

func (h *cookOrderUpdatedEventHandler) Handle(ctx context.Context, e *event.CookOrderUpdated) error {
	order, err := h.orderRepo.GetByID(ctx, e.OrderID)
	if err != nil {
		return errors.Wrap(err, "orderRepo.GetByID")
	}

	orderUp := event.OrderUp{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.Name,
		ItemType:   e.ItemType,
		TimeUp:     e.TimeUp,
		MadeBy:     e.MadeBy,
	}

	if err = order.Apply(&orderUp); err != nil {
		return errors.Wrap(err, "order.Apply")
	}

	_, err = h.orderRepo.Update(ctx, order)
	if err != nil {
		return errors.Wrap(err, "orderRepo.Update")
	}

	return nil
}

