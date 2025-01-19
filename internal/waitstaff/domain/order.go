
package domain

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	events "github.com/ntv97/atriaseniorliving/internal/pkg/event"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type Order struct {
	shared.AggregateRoot
	ID              uuid.UUID
	OrderTable      shared.OrderTable
	OrderName	string
	OrderStatus     shared.Status
	LineItems       []*LineItem
}

func NewOrder(
	orderTable shared.OrderTable,
	orderName string,
	orderStatus shared.Status,
) *Order {
	return &Order{
		ID:              uuid.New(),
		OrderTable:      orderTable,
		OrderName:       orderName,
		OrderStatus:     orderStatus,
	}
}

func CreateOrderFrom(
	ctx context.Context,
	request *PlaceOrderModel,
	itemsDomainSvc ItemsDomainService,
) (*Order, error) {
	order := NewOrder(request.OrderTable, request.OrderName.String(), shared.StatusInProcess)

	numberOfCookItems := len(request.CookItems) > 0
	numberOfChefItems := len(request.ChefItems) > 0
	numberOfWaitstaffItems := len(request.WaitstaffItems) > 0

	if numberOfCookItems {
		itemTypesRes, err := itemsDomainSvc.GetItemsByType(ctx, request, shared.Cook.String())
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.CookItems, func(item *OrderItemModel, _ int) {
			find, ok := lo.Find(itemTypesRes, func(i *ItemModel) bool {
				return i.ItemType == item.ItemType
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), find.OrderName.String(), shared.StatusInProcess, find.OrderType.String())

				event := events.CookOrdered{
					OrderID:    order.ID,
					ItemLineID: lineItem.ID,
					ItemType:   item.ItemType,
				}

				order.ApplyDomain(event)

				order.LineItems = append(order.LineItems, lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	if numberOfChefItems {
		itemTypesRes, err := itemsDomainSvc.GetItemsByType(ctx, request, shared.Chef.String())
		if err != nil {
			return nil, err
		}

		lo.ForEach(request.ChefItems, func(item *OrderItemModel, index int) {
			find, ok := lo.Find(itemTypesRes, func(i *ItemModel) bool {
				return i.ItemType == item.ItemType
			})

			if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), find.OrderName.String(), shared.StatusInProcess, find.OrderType.String())

				event := events.ChefOrdered{
					OrderID:    order.ID,
					ItemLineID: lineItem.ID,
					ItemType:   item.ItemType,
				}

				order.ApplyDomain(event)

				order.LineItems = append(order.LineItems, lineItem)
			}
		})

		if err != nil {
			return nil, err
		}
	}

	if numberOfWaitstaffItems {
                itemTypesRes, err := itemsDomainSvc.GetItemsByType(ctx, request, shared.Waitstaff.String())
                if err != nil {
                        return nil, err
                }

                lo.ForEach(request.WaitstaffItems, func(item *OrderItemModel, index int) {
                        find, ok := lo.Find(itemTypesRes, func(i *ItemModel) bool {
                                return i.ItemType == item.ItemType
                        })

                        if ok {
				lineItem := NewLineItem(item.ItemType, item.ItemType.String(), find.OrderName.String(), shared.StatusInProcess, find.OrderType.String())

                                event := events.WaitstaffOrdered{
                                        OrderID:    order.ID,
                                        ItemLineID: lineItem.ID,
                                        ItemType:   item.ItemType,
                                }

                                order.ApplyDomain(event)

                                order.LineItems = append(order.LineItems, lineItem)
                        }
                })

                if err != nil {
                        return nil, err
                }
        }

	return order, nil
}

func (o *Order) Apply(event *events.OrderUp) error {
	if len(o.LineItems) == 0 {
		return nil // we dont do anything
	}

	_, index, ok := lo.FindIndexOf(o.LineItems, func(i *LineItem) bool {
		return i.ItemType == event.ItemType
	})

	if !ok {
		return ErrItemNotFound
	}

	o.LineItems[index].ItemStatus = shared.StatusFulfilled

	if checkFulfilledStatus(o.LineItems) {
		o.OrderStatus = shared.StatusFulfilled
	}

	return nil
}

func checkFulfilledStatus(lineItems []*LineItem) bool {
	for _, item := range lineItems {
		if item.ItemStatus != shared.StatusFulfilled {
			return false
		}
	}

	return true
}

