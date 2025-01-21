
package domain

import (
	"context"
)

type (
	OrderRepo interface {
		Create(context.Context, *ChefOrder) error
	}
)

