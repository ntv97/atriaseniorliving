//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/ntv97/atriaseniorliving/cmd/items/config"
	"github.com/ntv97/atriaseniorliving/internal/items/app/router"
	"github.com/ntv97/atriaseniorliving/internal/items/infras/repo"
	 itemsUC "github.com/ntv97/atriaseniorliving/internal/items/usecases/items"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	grpcServer *grpc.Server,
) (*App, error) {
	panic(wire.Build(
		New,
		router.ItemGRPCServerSet,
		repo.RepositorySet,
		itemsUC.UseCaseSet,
	))
}
