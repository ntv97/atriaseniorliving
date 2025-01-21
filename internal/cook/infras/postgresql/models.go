
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package postgresql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CookCookOrder struct {
	ID       uuid.UUID    `json:"id"`
	OrderID  uuid.UUID    `json:"order_id"`
	ItemType int32        `json:"item_type"`
	ItemName string       `json:"item_name"`
	TimeUp   time.Time    `json:"time_up"`
	Created  time.Time    `json:"created"`
	Updated  sql.NullTime `json:"updated"`
}

