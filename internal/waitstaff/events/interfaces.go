package events

import (
	"context"

	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
)

type (
	CookOrderUpdatedEventHandler interface {
		Handle(context.Context, *event.CookOrderUpdated) error
	}

	ChefOrderUpdatedEventHandler interface {
		Handle(context.Context, *event.ChefOrderUpdated) error
	}
	WaitstaffOrderUpdatedEventHandler interface {
                Handle(context.Context, *event.WaitstaffOrderUpdated) error
        }
)
