package domain

import (
	"context"
)

type (
	ItemRepo interface {
		GetAll(context.Context) ([]*ItemTypeDto, error)
		GetByTypes(context.Context, []string) ([]*ItemDto, error)
	}
)
