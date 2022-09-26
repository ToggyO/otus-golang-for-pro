package entities

import (
	"database/sql"
	"time"
)

type BaseDbModel struct {
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
