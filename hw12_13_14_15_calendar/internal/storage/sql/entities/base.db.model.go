package entities

import (
	"database/sql"
	"time"
)

type BaseDBModel struct {
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
