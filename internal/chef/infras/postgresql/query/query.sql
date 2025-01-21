
-- name: CreateOrder :one

INSERT INTO
    chef.chef_orders (
        id,
        order_id,
        item_type,
        item_name,
        time_up,
        created,
        updated
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

