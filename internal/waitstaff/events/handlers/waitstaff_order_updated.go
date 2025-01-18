
package handlers

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/events"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
)

type waitstaffOrderUpdatedEventHandler struct {
	orderRepo orders.OrderRepo
}

var _ events.WaitstaffOrderUpdatedEventHandler = (*waitstaffOrderUpdatedEventHandler)(nil)

var WaitstaffOrderUpdatedEventHandlerSet = wire.NewSet(NewWaitstaffOrderUpdatedEventHandler)

func NewWaitstaffOrderUpdatedEventHandler(orderRepo orders.OrderRepo) events.WaitstaffOrderUpdatedEventHandler {
	return &waitstaffOrderUpdatedEventHandler{
		orderRepo: orderRepo,
	}
}

func (h *waitstaffOrderUpdatedEventHandler) Handle(ctx context.Context, e *event.WaitstaffOrderUpdated) error {
	order, err := h.orderRepo.GetByID(ctx, e.OrderID)
	if err != nil {
		return errors.Wrap(err, "orderRepo.GetOrderByID")
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

