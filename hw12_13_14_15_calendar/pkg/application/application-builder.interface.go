package application

import "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"

type IApplicationBuilder interface {
	SetServiceProvider(serviceProvider shared.IServiceProvider) IApplicationBuilder
	SetWebHost(host shared.IHost) IApplicationBuilder
	Build() (IApplication, error)
}
