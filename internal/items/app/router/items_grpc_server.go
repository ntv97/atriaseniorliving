package router

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/ntv97/atriaseniorliving/internal/items/usecases/items"
	"github.com/ntv97/atriaseniorliving/proto/gen"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ gen.ItemServiceServer = (*itemGRPCServer)(nil)

var ItemGRPCServerSet = wire.NewSet(NewItemGRPCServer)

type itemGRPCServer struct {
	gen.UnimplementedItemServiceServer
	uc items.UseCase
}

func NewItemGRPCServer(
	grpcServer *grpc.Server,
	uc items.UseCase,
) gen.ItemServiceServer {
	svc := itemGRPCServer{
		uc: uc,
	}

	gen.RegisterItemServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (g *itemGRPCServer) GetItemTypes(
	ctx context.Context,
	request *gen.GetItemTypesRequest,
) (*gen.GetItemTypesResponse, error) {
	slog.Info("gRPC client", "http_method", "GET", "http_name", "GetItemTypes")

	res := gen.GetItemTypesResponse{}

	results, err := g.uc.GetItemTypes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "itemGRPCServer-GetItemTypes")
	}

	for _, item := range results {
		res.ItemTypes = append(res.ItemTypes, &gen.ItemTypeDto{
			Name:  item.Name,
			Type:  int32(item.Type),
			Qty:   int32(item.Qty),
			Image: item.Image,
		})
	}

	return &res, nil
}

func (g *itemGRPCServer) GetItemsByType(
	ctx context.Context,
	request *gen.GetItemsByTypeRequest,
) (*gen.GetItemsByTypeResponse, error) {
	slog.Info("gRPC client", "http_method", "GET", "http_name", "GetItemsByType", "item_types", request.ItemTypes)

	res := gen.GetItemsByTypeResponse{}

	results, err := g.uc.GetItemsByType(ctx, request.ItemTypes)
	if err != nil {
		return nil, errors.Wrap(err, "itemGRPCServer-GetItemsByType")
	}

	for _, item := range results {
		res.Items = append(res.Items, &gen.ItemDto{
			Type:  int32(item.Type),
			Qty:   int32(item.Qty),
		})
	}

	return &res, nil
}
