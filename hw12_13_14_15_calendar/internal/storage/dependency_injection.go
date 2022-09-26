package storage

import (
	"context"

	sqlstorage "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/migrations"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

func AddStorage(
	ctx context.Context,
	serviceProvider shared.IServiceProvider,
	conf configuration.StorageConf,
) error {
	if conf.InMemory {
		return AddInMemoryStorage(ctx, serviceProvider)
	}

	return AddSqlStorage(ctx, serviceProvider, conf)
}

func AddSqlStorage(
	ctx context.Context,
	serviceProvider shared.IServiceProvider,
	conf configuration.StorageConf,
) error {
	var err error

	client := sqlstorage.NewDbClient(conf)
	if err = client.Connect(ctx); err != nil {
		return err
	}

	bindDbClient := func() shared.IDbClient {
		return client
	}

	err = serviceProvider.AddService(&shared.ServiceDescriptor{Service: bindDbClient})
	err = serviceProvider.AddService(&shared.ServiceDescriptor{Service: sqlstorage.NewEventsRepository})
	if err != nil {
		return err
	}

	return err
}

func AddInMemoryStorage(
	_ context.Context,
	serviceProvider shared.IServiceProvider,
) error {
	err := serviceProvider.AddService(&shared.ServiceDescriptor{Service: sqlstorage.NewEventsRepository})
	return err
}

func AddMigrationRunner(
	_ context.Context,
	serviceProvider shared.IServiceProvider,
	conf configuration.StorageConf,
) error {
	return serviceProvider.AddService(&shared.ServiceDescriptor{
		Service: func() migrations.IMigrationRunner {
			return migrations.NewMigrationRunner(conf.Dialect, shared.CreatePgConnectionString(conf))
		},
	})
}
