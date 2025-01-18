
package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/domain"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/infras/postgresql"
	"github.com/ntv97/atriaseniorliving/internal/waitstaff/usecases/orders"
	shared "github.com/ntv97/atriaseniorliving/internal/pkg/shared_kernel"
	"github.com/ntv97/atriaseniorliving/pkg/postgres"
)

const _defaultEntityCap = 64

type orderRepo struct {
	pg postgres.DBEngine
}

var _ orders.OrderRepo = (*orderRepo)(nil)

var RepositorySet = wire.NewSet(NewOrderRepo)

func NewOrderRepo(pg postgres.DBEngine) orders.OrderRepo {
	return &orderRepo{pg: pg}
}

func (d *orderRepo) GetAll(ctx context.Context) ([]*domain.Order, error) {
	querier := postgresql.New(d.pg.GetDB())

	results, err := querier.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetAll")
	}

	uniqueResults := lo.UniqBy(results, func(x postgresql.GetAllRow) string {
		return x.ID.String()
	})
	orders := lo.Map(uniqueResults, func(x postgresql.GetAllRow, _ int) *domain.Order {
		return &domain.Order{
                        ID:              x.ID,
                        OrderTable:      shared.OrderTable(x.OrderTable),
                        OrderName:       shared.OrderName(x.OrderName),
                        OrderStatus:     shared.Status(x.OrderStatus),
                }
	})
	lineItems := lo.Map(results, func(x postgresql.GetAllRow, _ int) *domain.LineItem {
		return &domain.LineItem{
                        ID:             x.LineItemID.UUID,
                        ItemType:       shared.ItemType(x.ItemType),
                        ItemName:       x.ItemName,
                        OrderName:      x.OrderName,
                        ItemStatus:     shared.Status(x.ItemStatus),
                        OrderType:      x.OrderType,
                        OrderID:        x.ID,
                }
	})
	entities := make([]*domain.Order, 0, _defaultEntityCap)

	for _, o := range orders {
		order := &domain.Order{
                        ID:              o.ID,
                        OrderTable:      o.OrderTable,
                        OrderName:       o.OrderName,
                        OrderStatus:     o.OrderStatus,
                }

		filters := lo.Filter(lineItems, func(x *domain.LineItem, _ int) bool {
			return x.OrderID.ID() == o.ID.ID()
		})

		for _, ol := range filters {
			order.LineItems = append(order.LineItems, &domain.LineItem{
				ID:             ol.ID,
				ItemType:       ol.ItemType,
				ItemName:       ol.ItemName,
				OrderName:      ol.OrderName,
				ItemStatus:     ol.ItemStatus,
				OrderType:      ol.OrderType,
				OrderID:        ol.OrderID,
			})
		}

		entities = append(entities, order)
	}

	return entities, nil
}

func (d *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	querier := postgresql.New(d.pg.GetDB())

	results, err := querier.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetAll")
	}

	uniqueResults := lo.UniqBy(results, func(x postgresql.GetByIDRow) string {
		return x.ID.String()
	})

	orders := lo.Map(uniqueResults, func(x postgresql.GetByIDRow, _ int) *domain.Order {
		return &domain.Order{
                        ID:              x.ID,
                        OrderTable:      shared.OrderTable(x.OrderTable),
                        OrderName:       shared.OrderName(x.OrderName),
                        OrderStatus:     shared.Status(x.OrderStatus),
		}
	})
	lineItems := lo.Map(results, func(x postgresql.GetByIDRow, _ int) *domain.LineItem {
		return &domain.LineItem{
                        ID:             x.LineItemID.UUID,
                        ItemType:       shared.ItemType(x.ItemType),
                        ItemName:       x.ItemName,
                        OrderName:      x.OrderName,
                        ItemStatus:     shared.Status(x.ItemStatus),
                        OrderType:      x.OrderType,
                        OrderID:        x.ID,
                }
	})

	if len(orders) == 0 {
		return nil, nil
	}

	order := &domain.Order{
		ID:              orders[0].ID,
		OrderTable:      orders[0].OrderTable,
		OrderName:       orders[0].OrderName,
		OrderStatus:     orders[0].OrderStatus,
	}

	for _, ol := range lineItems {
		order.LineItems = append(order.LineItems, &domain.LineItem{
			ID:             ol.ID,
			ItemType:       ol.ItemType,
			ItemName:       ol.ItemName,
			OrderName:      ol.OrderName,
			ItemStatus:     ol.ItemStatus,
			OrderType:      ol.OrderType,
			OrderID:        ol.OrderID,
		})
	}

	return order, nil
}

func (d *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	db := d.pg.GetDB()
	querier := postgresql.New(db)

	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "coolOrderedEventHandler.Handle")
	}

	qtx := querier.WithTx(tx)

	_, err = qtx.CreateOrder(ctx, postgresql.CreateOrderParams{
		ID:              order.ID,
                OrderTable:     int32(order.OrderTable),
                OrderName:      order.OrderName,
                OrderStatus:     int32(order.OrderStatus),
                Updated: sql.NullTime{
                        Time:  time.Now(),
                        Valid: true,
                },
	})
	if err != nil {
		return errors.Wrap(err, "qtx.CreateOrder(ctx, postgresql.CreateOrderParams{})")
	}

	// continue to insert order items
	for _, item := range order.LineItems {
		_, err = qtx.InsertItemLine(ctx, postgresql.InsertItemLineParams{
			ID:              item.ID,
			ItemType: 	int32(item.ItemType),
			ItemName:	item.ItemName,
			OrderName:	item.OrderName,
			ItemStatus:     int32(item.ItemStatus),
			OrderType:	item.OrderType,
			OrderID: uuid.NullUUID{
                                UUID:  order.ID,
                                Valid: true,
                        },
			Created: time.Now(),
                        Updated: sql.NullTime{
                                Time:  time.Now(),
                                Valid: true,
                        },
		})

		if err != nil {
			return errors.Wrap(err, "qtx.InsertItemLine(ctx, postgresql.InsertItemLineParams{})")
		}
	}

	return tx.Commit()
}

func (d *orderRepo) Update(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	db := d.pg.GetDB()
	querier := postgresql.New(db)

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "cookOrderedEventHandler.Handle")
	}

	qtx := querier.WithTx(tx)

	err = qtx.UpdateOrder(ctx, postgresql.UpdateOrderParams{
		ID:          order.ID,
		OrderStatus: int32(order.OrderStatus),
		Updated: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.UpdateOrder(ctx, postgresql.UpdateOrderParams{})")
	}

	// continue to insert order items
	for _, item := range order.LineItems {
		err = qtx.UpdateItemLine(ctx, postgresql.UpdateItemLineParams{
			ID:         item.ID,
			ItemStatus: int32(item.ItemStatus),
			Updated: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})

		if err != nil {
			return nil, errors.Wrap(err, "qtx.UpdateItemLine(ctx, postgresql.UpdateItemLineParams{})")
		}
	}

	return nil, tx.Commit()
}

