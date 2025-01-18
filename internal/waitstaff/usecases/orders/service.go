
package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	"golang.org/x/exp/slog"
)

type usecase struct {
	orderRepo        OrderRepo
	itemsDomainSvc domain.ItemsDomainService
	cookEventPub  CookEventPublisher
	chefEventPub  ChefEventPublisher
	waitstaffEventPub WaitstaffEventPublisher
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	orderRepo OrderRepo,
	itemsDomainSvc domain.ItemsDomainService,
	cookEventPub CookEventPublisher,
	chefEventPub ChefEventPublisher,
	waistaffEventPub WaitstaffEventPublisher,
) UseCase {
	return &usecase{
		orderRepo:        orderRepo,
		itemsDomainSvc: itemsDomainSvc,
		cookEventPub:  cookEventPub,
		chefEventPub:  chefEventPub,
		waitstaffEventPub: waitstaffEventPub,
	}
}

func (uc *usecase) GetListOrderFulfillment(ctx context.Context) ([]*domain.Order, error) {
	entities, err := uc.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("orderRepo.GetAll: %w", err)
	}

	return entities, nil
}

func (uc *usecase) PlaceOrder(ctx context.Context, model *domain.PlaceOrderModel) error {
	order, err := domain.CreateOrderFrom(ctx, model, uc.itemsDomainSvc)
	if err != nil {
		return errors.Wrap(err, "domain.CreateOrderFrom")
	}

	err = uc.orderRepo.Create(ctx, order)
	if err != nil {
		return errors.Wrap(err, "orderRepo.Create")
	}

	slog.Debug("order created", "order", *order)

	// todo: it might cause dual-write problem, but we accept it temporary
	for _, event := range order.DomainEvents() {
		if event.Identity() == "CookOrdered" {
			eventBytes, err := json.Marshal(event)
			if err != nil {
				return errors.Wrap(err, "json.Marshal[event]")
			}

			uc.cookEventPub.Publish(ctx, eventBytes, "text/plain")
		}

		if event.Identity() == "ChefOrdered" {
			eventBytes, err := json.Marshal(event)
			if err != nil {
				return errors.Wrap(err, "json.Marshal[event]")
			}

			uc.chefEventPub.Publish(ctx, eventBytes, "text/plain")
		}
		if event.Identity() == "WaitstaffOrdered" {
                        eventBytes, err := json.Marshal(event)
                        if err != nil {
                                return errors.Wrap(err, "json.Marshal[event]")
                        }

                        uc.waitstaffEventPubEventPub.Publish(ctx, eventBytes, "text/plain")
                }
	}

	return nil
}

