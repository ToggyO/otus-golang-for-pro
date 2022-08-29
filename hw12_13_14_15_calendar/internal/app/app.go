package app

import (
	appcore "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/app/core"
	sqlstorage "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/repositories"
)

type Application struct { // TODO
	serviceProvider *serviceProvider
	logger          appcore.ILogger
	storage         appcore.IStorage
}

func NewApplication(logger appcore.ILogger, , inMemoryStorage bool) *Application {
	sqlStorageModule := sqlstorage.NewSqlStorageModule()
	sp := &serviceProvider{
		db: storage,
		eventsRepository: repositories.NewEventsRepository(,
	}
	return &Application{}
}

func (a *Application) prepareServiceProvider() *serviceProvider {
	return
}

//func (a *App) CreateEvent(ctx context.Context, id, title string) error {
//	// TODO
//	return nil
//	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
//}

// TODO
