
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one

INSERT INTO
    "order".orders (
        id,
	order_table,
	order_name,
	order_type,
        order_status,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, order_table, order_name, order_type, order_status, updated
`

type CreateOrderParams struct {
	ID              uuid.UUID    `json:"id"`
        OrderTable      int32        `json:"order_table"`
        OrderName       string       `json:"order_name"`
        OrderType       string       `json:"order_type"`
        OrderStatus     int32        `json:"order_status"`
        Updated         sql.NullTime `json:"updated"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (OrderOrder, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.ID,
		arg.OrderTable,
		arg.OrderName,
		arg.OrderType,
		arg.OrderStatus,
		arg.Updated,
	)
	var i OrderOrder
	err := row.Scan(
		&i.ID,
		&i.OrderTable,
		&i.OrderName,
		&i.OrderType,
		&i.OrderStatus,
		&i.Updated,
	)
	return i, err
}

const getAll = `-- name: GetAll :many

SELECT
    o.id,
    order_table,
    order_name,
    order_status,
    l.id as "line_item_id",
    item_type,
    item_name,
    order_name,
    item_status,
    order_type
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id
`

type GetAllRow struct {
	ID              uuid.UUID     `json:"id"`
	OrderTable      int32         `json:"order_table"`
	OrderName       string        `json:"order_name"`
        OrderStatus     int32         `json:"order_status"`
	LineItemID      uuid.NullUUID `json:"line_item_id"`
	ItemType        int32         `json:"item_type"`
        ItemName        string        `json:"item_name"`
        OrderName       string        `json:"order_name"`
        ItemStatus      int32         `json:"item_status"`
	OrderType       string         `json:"order_type"`
}

func (q *Queries) GetAll(ctx context.Context) ([]GetAllRow, error) {
	rows, err := q.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllRow
	for rows.Next() {
		var i GetAllRow
		if err := rows.Scan(
			&i.ID,
			&i.OrderTable,
			&i.OrderName,
			&i.OrderStatus,
			&i.LineItemID,
			&i.ItemType,
			&i.ItemName,
			&i.OrderName,
			&i.ItemStatus,
			&i.OrderType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getByID = `-- name: GetByID :many

SELECT
    o.id,
    order_table,
    order_name,
    order_status,
    l.id as "line_item_id",
    item_type,
    item_name,
    order_name,
    item_status,
    order_type
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id
WHERE o.id = $1
`

type GetByIDRow struct {
	ID              uuid.UUID     `json:"id"`
        OrderTable      int32         `json:"order_table"`
        OrderName       string        `json:"order_name"`
        OrderStatus     int32         `json:"order_status"`
        LineItemID      uuid.NullUUID `json:"line_item_id"`
        ItemType        int32         `json:"item_type"`
        ItemName        string        `json:"item_name"`
        OrderName       string        `json:"order_name"`
        ItemStatus      int32         `json:"item_status"`
        OrderType       string         `json:"order_type"`
}

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) ([]GetByIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetByIDRow
	for rows.Next() {
		var i GetByIDRow
		if err := rows.Scan(
			&i.ID,
                        &i.OrderTable,
                        &i.OrderName,
                        &i.OrderStatus,
                        &i.LineItemID,
                        &i.ItemType,
                        &i.ItemName,
                        &i.OrderName,
                        &i.ItemStatus,
                        &i.OrderType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertItemLine = `-- name: InsertItemLine :one

INSERT INTO
    "order".line_items (
        id,
        item_type,
        item_name,
        order_name,
        item_status,
        order_type,
        order_id,
        created,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, item_type, item_name, order_name, item_status, order_type, order_id, created, updated
`

type InsertItemLineParams struct {
	ID             uuid.UUID     `json:"id"`
        ItemType       int32         `json:"item_type"`
        ItemName       string        `json:"item_name"`
        OrderName      string        `json:"order_name"`
        ItemStatus     int32         `json:"item_status"`
        OrderType      string         `json:"order_type"`
        OrderID        uuid.NullUUID `json:"order_id"`
        Created        time.Time     `json:"created"`
        Updated        sql.NullTime  `json:"updated"`
}

func (q *Queries) InsertItemLine(ctx context.Context, arg InsertItemLineParams) (OrderLineItem, error) {
	row := q.db.QueryRowContext(ctx, insertItemLine,
		arg.ID,
		arg.ItemType,
		arg.ItemName,
		arg.OrderName,
		arg.ItemStatus,
		arg.OrderType,
		arg.OrderID,
		arg.Created,
		arg.Updated,
	)
	var i OrderLineItem
	err := row.Scan(
		&i.ID,
		&i.ItemType,
		&i.ItemName,
		&i.OrderName,
		&i.ItemStatus,
		&i.OrderType,
		&i.OrderID,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const updateItemLine = `-- name: UpdateItemLine :exec

UPDATE "order".line_items
SET
    item_status = $2,
    updated = $3
WHERE id = $1
`

type UpdateItemLineParams struct {
	ID         uuid.UUID    `json:"id"`
	ItemStatus int32        `json:"item_status"`
	Updated    sql.NullTime `json:"updated"`
}

func (q *Queries) UpdateItemLine(ctx context.Context, arg UpdateItemLineParams) error {
	_, err := q.db.ExecContext(ctx, updateItemLine, arg.ID, arg.ItemStatus, arg.Updated)
	return err
}

const updateOrder = `-- name: UpdateOrder :exec

UPDATE "order".orders
SET
    order_status = $2,
    updated = $3
WHERE id = $1
`

type UpdateOrderParams struct {
	ID          uuid.UUID    `json:"id"`
	OrderStatus int32        `json:"order_status"`
	Updated     sql.NullTime `json:"updated"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) error {
	_, err := q.db.ExecContext(ctx, updateOrder, arg.ID, arg.OrderStatus, arg.Updated)
	return err
}

