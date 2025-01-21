package eventhandlers

import (
	"context"

	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
)

type CookOrderedEventHandler interface {
	Handle(context.Context, event.CookOrdered) error
}
