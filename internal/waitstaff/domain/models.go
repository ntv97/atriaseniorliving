
package domain

import (
	"time"

	"github.com/google/uuid"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type PlaceOrderModel struct {
	CommandType     shared.CommandType
	OrderTable      shared.OrderTable
	OrderName       shared.OrderName
	CookItems    	[]*OrderItemModel
	ChefItems    	[]*OrderItemModel
	WaitstaffITems  []*OrderItemModel
	Timestamp       time.Time
}

type OrderItemModel struct {
	ItemType shared.ItemType
}

type ItemModel struct {
	ItemType 	shared.ItemType
	OrderTable 	shared.OrderTable
	OrderName 	shared.OrderName
	Qty		int32
}

