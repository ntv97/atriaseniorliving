
-- name: GetAll :many

SELECT
    o.id,
    order_table,
    order_name,
    order_status,
    l.id as "line_item_id",
    item_type,
    item_name,
    item_order_name,
    item_status,
    order_type
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id;

-- name: GetByID :many

SELECT
    o.id,
    order_table,
    order_name,
    order_status,
    l.id as "line_item_id",
    item_type,
    item_name,
    item_order_name,
    item_status,
    order_type
FROM "order".orders o
    LEFT JOIN "order".line_items l ON o.id = l.order_id
WHERE o.id = $1;

-- name: CreateOrder :one

INSERT INTO
    "order".orders (
        id,
        order_table,
        order_name,
        order_status,
        updated
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: InsertItemLine :one

INSERT INTO
    "order".line_items (
        id,
        item_type,
        item_name,
        item_order_name,
        item_status,
        order_type,
        order_id,
        created,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateOrder :exec

UPDATE "order".orders
SET
    order_status = $2,
    updated = $3
WHERE id = $1;

-- name: UpdateItemLine :exec

UPDATE "order".line_items
SET
    item_status = $2,
    updated = $3
WHERE id = $1;

