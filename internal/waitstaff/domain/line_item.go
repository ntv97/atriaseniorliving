
package domain

import (
	"github.com/google/uuid"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type LineItem struct {
	ID             uuid.UUID
        ItemType       shared.ItemType        
        ItemName       string        
        OrderName      string        
        ItemStatus     shared.Status       
        OrderType      string         
        OrderID        uuid.UUID
}

func NewLineItem(itemType shared.ItemType, itemName string, orderName string, itemStatus shared.Status, orderType string) *LineItem {
	return &LineItem{
		ID:             uuid.New(),
		ItemType:       itemType,
		ItemName:       itemName,
	        OrderName:      orderName,
		ItemStatus:     itemStatus,
		OrderType:      orderType,
	}
}

