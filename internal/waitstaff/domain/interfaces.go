
package domain

import (
	"context"
)

type (
	ItemsDomainService interface {
		GetItemsByType(context.Context, *PlaceOrderModel, string) ([]*ItemModel, error)
	}
)

