package main

import (
	"context"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/api"
	infrastructure "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/infrastructure/services"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/migrations"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type Startup struct{}

func (s Startup) ConfigureServices(
	ctx context.Context,
	configuration configuration.Configuration,
	serviceProvider shared.IServiceProvider,
) error {
	err := infrastructure.AddInfrastructureServices(ctx, serviceProvider, configuration)
	err = storage.AddStorage(ctx, serviceProvider, configuration.Storage)
	err = storage.AddMigrationRunner(ctx, serviceProvider, configuration.Storage)
	err = api.AddHttpHandler(ctx, serviceProvider)
	return err
}

func (s Startup) AfterApplicationStartup(_ context.Context, serviceProvider shared.IServiceProvider) error {
	var migrationRunner migrations.IMigrationRunner
	err := serviceProvider.GetService(func(mr migrations.IMigrationRunner) {
		migrationRunner = mr
	})
	if err != nil {
		return err
	}

	err = migrationRunner.MigrateUp(context.Background())
	return err
}

func (s Startup) BeforeApplicationShutdown(ctx context.Context, serviceProvider shared.IServiceProvider) error {
	var client shared.IDbClient
	err := serviceProvider.GetService(func(c shared.IDbClient) {
		client = c
	})

	if err != nil {
		return err
	}

	err = client.Close(ctx)
	return err
}
