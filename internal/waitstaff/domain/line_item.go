
package domain

import (
	"github.com/google/uuid"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
)

type LineItem struct {
	ID             uuid.UUID
        ItemType       int32         
        ItemName       string        
        OrderName      string        
        ItemStatus     int32        
        OrderType      string         
        OrderID        uuid.NullUUID 
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

