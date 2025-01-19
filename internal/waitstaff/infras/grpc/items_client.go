
package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/ntv97/atriaseniorliving/cmd/waitstaff/config"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
	gen "github.com/ntv97/atriaseniorliving/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type itemsGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.ItemsDomainService = (*itemsGRPCClient)(nil)

var ItemsGRPCClientSet = wire.NewSet(NewGRPCItemsClient)

func NewGRPCItemsClient(cfg *config.Config) (domain.ItemsDomainService, error) {
	conn, err := grpc.Dial(cfg.ItemsClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &itemsGRPCClient{
		conn: conn,
	}, nil
}

func (p *itemsGRPCClient) GetItemsByType(
	ctx context.Context,
	model *domain.PlaceOrderModel,
	orderType string,
) ([]*domain.ItemModel, error) {
	c := gen.NewItemServiceClient(p.conn)

	itemTypes := ""
	if orderType == "Cook" {
		itemTypes = lo.Reduce(model.CookItems, func(agg string, item *domain.OrderItemModel, _ int) string {
			return fmt.Sprintf("%s,%s", agg, item.ItemType.String())
		}, "")
	} else if orderType == "Chef" {
		itemTypes = lo.Reduce(model.ChefItems, func(agg string, item *domain.OrderItemModel, _ int) string {
			return fmt.Sprintf("%s,%s", agg, item.ItemType.String())
		}, "")
	} else {
                itemTypes = lo.Reduce(model.WaitstaffItems, func(agg string, item *domain.OrderItemModel, _ int) string {
                        return fmt.Sprintf("%s,%s", agg, item.ItemType.String())
                }, "")
        }

	res, err := c.GetItemsByType(ctx, &gen.GetItemsByTypeRequest{ItemTypes: strings.TrimLeft(itemTypes, ",")})
	if err != nil {
		return nil, errors.Wrap(err, "itemsGRPCClient-c.GetItemsByType")
	}

	results := make([]*domain.ItemModel, 0)
	for _, item := range res.Items {
		results = append(results, &domain.ItemModel{
			ItemType: shared.ItemType(item.Type),
			Qty:    item.Qty,
		})
	}

	return results, nil
}

