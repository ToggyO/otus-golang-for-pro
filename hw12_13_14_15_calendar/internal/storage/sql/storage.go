package sqlstorage

import (
	"context"
	"fmt"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/jmoiron/sqlx"
	"strings"
)

type PgDbClient struct { // TODO
	host     string
	port     int
	user     string
	password string
	dbName   string

	db *sqlx.DB
}

func NewDbClient(conf config.StorageConf) *PgDbClient {
	return &PgDbClient{
		host:     conf.Host,
		port:     conf.Port,
		user:     conf.User,
		password: conf.Password,
		dbName:   conf.DbName,
	}
}

func (pg *PgDbClient) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", pg.connectionString())
	if err != nil {
		return nil
	}

	pg.db = db
	return nil
}

func (pg *PgDbClient) Close(_ context.Context) error {
	// TODO: check
	return pg.db.Close()
}

func (pg *PgDbClient) GetConnection() *sqlx.DB {
	return pg.db
}

func (pg *PgDbClient) connectionString() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("host=%s port=%d", pg.host, pg.port))

	if pg.user != "" {
		sb.WriteString(fmt.Sprintf(" user=%s", pg.user))
	}
	if pg.password != "" {
		sb.WriteString(fmt.Sprintf(" password=%s", pg.user))
	}
	if pg.dbName != "" {
		sb.WriteString(fmt.Sprintf(" dbname=%s", pg.dbName))
	}

	return sb.String()
}
