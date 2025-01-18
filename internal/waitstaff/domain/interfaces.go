
package domain

import (
	"context"
)

type (
	ItemsDomainService interface {
		GetItemsByType(context.Context, *PlaceOrderModel, bool) ([]*ItemModel, error)
	}
)

