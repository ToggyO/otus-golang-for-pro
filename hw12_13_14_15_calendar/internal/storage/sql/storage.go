package sqlstorage

import (
	"context"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // lint:ignore revive
)

type PgDBClient struct {
	dialect          string
	connectionString string

	db *sqlx.DB
}

func NewDBClient(conf configuration.StorageConf) shared.IDbClient {
	return &PgDBClient{
		dialect:          conf.Dialect,
		connectionString: shared.CreatePgConnectionString(conf),
	}
}

func (pg *PgDBClient) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(
		ctx,
		pg.dialect,
		pg.connectionString)
	if err != nil {
		return err
	}

	pg.db = db
	return nil
}

func (pg *PgDBClient) Close(_ context.Context) error {
	return pg.db.Close()
}

func (pg *PgDBClient) GetConnection() interface{} {
	return pg.db
}
