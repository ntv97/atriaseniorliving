
package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type ChefOrder struct {
	shared.AggregateRoot
	ID       uuid.UUID
	OrderID  uuid.UUID
	ItemName string
	ItemType shared.ItemType
	TimeUp   time.Time
	Created  time.Time
	Updated  time.Time
}

func NewChefOrder(e event.ChefOrdered) ChefOrder {
	timeIn := time.Now()

	delay := calculateDelay(e.ItemType)
	time.Sleep(delay) // simulate the delay when makes the drink

	timeUp := time.Now().Add(delay)

	order := ChefOrder{
		ID:       e.ItemLineID,
		OrderID:  e.OrderID,
		ItemName: e.ItemType.String(),
		ItemType: e.ItemType,
		TimeUp:   timeUp,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	orderUpdatedEvent := event.ChefOrderUpdated{
		OrderID:    e.OrderID,
		ItemLineID: e.ItemLineID,
		Name:       e.ItemType.String(),
		ItemType:   e.ItemType,
		MadeBy:     "teesee",
		TimeIn:     timeIn,
		TimeUp:     timeUp,
	}

	order.ApplyDomain(&orderUpdatedEvent)

	return order
}

func calculateDelay(itemType shared.ItemType) time.Duration {
	switch itemType {
	case shared.ItemTypeFourCheeseRavioli:
		return 7 * time.Second
	case shared.ItemTypeOmelet:
		return 7 * time.Second
	default:
		return 3 * time.Second
	}
}

