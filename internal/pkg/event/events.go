package event

import (
	"time"

	"github.com/google/uuid"

	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type CookOrdered struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	ItemType   shared.ItemType `json:"itemType"`
}

func (e CookOrdered) Identity() string {
	return "CookOrdered"
}

type ChefOrdered struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	ItemType   shared.ItemType `json:"itemType"`
}

func (e ChefOrdered) Identity() string {
	return "ChefOrdered"
}

type WaitstaffOrdered struct {
        shared.DomainEvent
        OrderID    uuid.UUID       `json:"orderId"`
        ItemLineID uuid.UUID       `json:"itemLineId"`
        ItemType   shared.ItemType `json:"itemType"`
}

func (e WaitstaffOrdered) Identity() string {
        return "WaitstaffOrdered"
}

type CookOrderUpdated struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeIn     time.Time       `json:"timeIn"`
	MadeBy     string          `json:"madeBy"`
	TimeUp     time.Time       `json:"timeUp"`
}

func (e *CookOrderUpdated) Identity() string {
	return "CookOrderUpdated"
}

type ChefOrderUpdated struct {
	shared.DomainEvent
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeIn     time.Time       `json:"timeIn"`
	MadeBy     string          `json:"madeBy"`
	TimeUp     time.Time       `json:"timeUp"`
}

func (e *ChefOrderUpdated) Identity() string {
	return "ChefOrderUpdated"
}

type WaitstaffOrderUpdated struct {
        shared.DomainEvent
        OrderID    uuid.UUID       `json:"orderId"`
        ItemLineID uuid.UUID       `json:"itemLineId"`
        Name       string          `json:"name"`
        ItemType   shared.ItemType `json:"itemType"`
        TimeIn     time.Time       `json:"timeIn"`
        MadeBy     string          `json:"madeBy"`
        TimeUp     time.Time       `json:"timeUp"`
}

func (e WaitstaffOrderUpdated) Identity() string {
        return "WaitstaffOrderUpdated"
}

type OrderUp struct {
	OrderID    uuid.UUID       `json:"orderId"`
	ItemLineID uuid.UUID       `json:"itemLineId"`
	Name       string          `json:"name"`
	ItemType   shared.ItemType `json:"itemType"`
	TimeUp     time.Time       `json:"timeUp"`
	MadeBy     string          `json:"madeBy"`
}

func (e *OrderUp) Identity() string {
	return "OrderUp"
}
