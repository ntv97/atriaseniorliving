
package eventhandlers

import (
	"context"

	"github.com/ntv97/atriaseniorliving/internal/pkg/event"
)

type ChefOrderedEventHandler interface {
	Handle(context.Context, event.ChefOrdered) error
}

