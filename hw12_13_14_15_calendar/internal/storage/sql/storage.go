package sqlstorage

import (
	"context"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
)

type PgDbClient struct {
	dialect          string
	connectionString string

	db *sqlx.DB
}

func NewDbClient(conf configuration.StorageConf) shared.IDbClient {
	return &PgDbClient{
		dialect:          conf.Dialect,
		connectionString: shared.CreatePgConnectionString(conf),
	}
}

func (pg *PgDbClient) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(
		ctx,
		pg.dialect,
		pg.connectionString)
	if err != nil {
		return nil
	}

	pg.db = db
	return nil
}

func (pg *PgDbClient) Close(_ context.Context) error {
	return pg.db.Close()
}

func (pg *PgDbClient) GetConnection() interface{} {
	return pg.db
}
