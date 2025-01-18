package app

import (
	"github.com/ntv97/atriaseniorliving/cmd/items/config"
	itemUC "github.com/ntv97/atriaseniorliving/internal/items/usecases/items"
	"github.com/ntv97/atriaseniorliving/proto/gen"
)

type App struct {
	Cfg               *config.Config
	UC                itemUC.UseCase
	ItemGRPCServer    gen.ItemServiceServer
}

func New(
	cfg *config.Config,
	uc itemUC.UseCase,
	itemGRPCServer gen.ItemServiceServer,
) *App {
	return &App{
		Cfg:               cfg,
		UC:                uc,
		ItemGRPCServer:    itemGRPCServer,
	}
}
