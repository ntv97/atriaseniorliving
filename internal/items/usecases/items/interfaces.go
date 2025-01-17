package items

import (
	"context"

	"github.com/ntv97/atriaseniorliving/internal/items/domain"
)

type UseCase interface {
	GetItemTypes(context.Context) ([]*domain.ItemTypeDto, error)
	GetItemsByType(context.Context, string) ([]*domain.ItemDto, error)
}
