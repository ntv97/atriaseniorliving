
package router

import (
	"context"
	"fmt"

	//"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/ntv97/atriaseniorliving/cmd/waitstaff/config"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
	gen "github.com/ntv97/atriaseniorliving/proto/gen"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type waitstaffGRPCServer struct {
	gen.UnimplementedWaitstaffServiceServer
	cfg *config.Config
	uc  orders.UseCase
}

var _ gen.WaitstaffServiceServer = (*waitstaffGRPCServer)(nil)

var WaitstaffGRPCServerSet = wire.NewSet(NewGRPCWaitstaffServer)

func NewGRPCWaitstaffServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc orders.UseCase,
) gen.WaitstaffServiceServer {
	svc := waitstaffGRPCServer{
		cfg: cfg,
		uc:  uc,
	}

	gen.RegisterWaitstaffServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (g *waitstaffGRPCServer) GetListOrderFulfillment(
	ctx context.Context,
	request *gen.GetListOrderFulfillmentRequest,
) (*gen.GetListOrderFulfillmentResponse, error) {
	slog.Info("GET: GetListOrderFulfillment")

	res := gen.GetListOrderFulfillmentResponse{}

	entities, err := g.uc.GetListOrderFulfillment(ctx)
	if err != nil {
		return nil, fmt.Errorf("uc.GetListOrderFulfillment: %w", err)
	}

	for _, entity := range entities {
		res.Orders = append(res.Orders, &gen.OrderDto{
			Id:              entity.ID.String(),
			OrderTable:     int32(entity.OrderTable),
			OrderName:     entity.OrderName,
			LineItems: lo.Map(entity.LineItems, func(item *domain.LineItem, _ int) *gen.LineItemDto {
				return &gen.LineItemDto{
					Id:             item.ID.String(),
					ItemType:       int32(item.ItemType),
					ItemName:       item.ItemName,
					OrderName:      item.OrderName,
					ItemStatus:     int32(item.ItemStatus),
					OrderType:      item.OrderType,
				}
			}),
		})
	}

	return &res, nil
}

func (g *waitstaffGRPCServer) PlaceOrder(
	ctx context.Context,
	request *gen.PlaceOrderRequest,
) (*gen.PlaceOrderResponse, error) {
	slog.Info("POST: PlaceOrder")

	model := domain.PlaceOrderModel{
		CommandType:     shared.CommandType(request.CommandType),
		OrderTable:     shared.OrderTable(request.OrderTable),
		OrderName:        shared.OrderName(request.OrderName),
		Timestamp:       request.Timestamp.AsTime(),
	}

	for _, cook := range request.CookItems {
		model.CookItems = append(model.CookItems, &domain.OrderItemModel{
			ItemType: shared.ItemType(cook.ItemType),
		})
	}

	for _, chef := range request.ChefItems {
		model.ChefItems = append(model.ChefItems, &domain.OrderItemModel{
			ItemType: shared.ItemType(chef.ItemType),
		})
	}

	for _, waitstaff := range request.WaitstaffItems {
                model.ChefItems = append(model.WaitstaffItems, &domain.OrderItemModel{
                        ItemType: shared.ItemType(waitstaff.ItemType),
                })
        }

	err := g.uc.PlaceOrder(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.PlaceOrder")
	}

	res := gen.PlaceOrderResponse{}

	return &res, nil
}

