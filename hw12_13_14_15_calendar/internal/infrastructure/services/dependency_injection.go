package infrastructure

import (
	"context"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/infrastructure/services/logger"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

func AddInfrastructureServices(
	_ context.Context,
	serviceProvider shared.IServiceProvider,
	configuration configuration.Configuration,
) error {
	bindLogger := func() shared.ILogger {
		return logger.NewLogger(configuration.Logger.Level, configuration.IsDev, nil)
	}

	err := serviceProvider.AddService(&shared.ServiceDescriptor{Service: bindLogger})

	return err
}
