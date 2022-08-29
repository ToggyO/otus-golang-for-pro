package sqlstorage

import (
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/config"
	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/repositories"
)

type Module struct {
	client *PgDbClient

	eventsRepository domain.IEventsRepository
}

func NewSqlStorageModule(conf config.StorageConf) *Module {
	client := NewDbClient(conf)
	conn := client.GetConnection()
	return &Module{
		client: client,

		eventsRepository: repositories.NewEventsRepository(conn),
	}
}
