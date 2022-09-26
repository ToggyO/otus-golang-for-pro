package application

import (
	"context"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type IStartup interface {
	ConfigureServices(ctx context.Context, configuration configuration.Configuration, serviceProvider shared.IServiceProvider) error
	AfterApplicationStartup(ctx context.Context, serviceProvider shared.IServiceProvider) error
	BeforeApplicationShutdown(ctx context.Context, serviceProvider shared.IServiceProvider) error
}
