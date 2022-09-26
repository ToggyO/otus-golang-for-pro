package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/configuration"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/ioc"
	httpserver "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/server/http"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type Application struct {
	startup       IStartup
	host          shared.IHost
	configuration configuration.Configuration
	container     shared.IServiceProvider
	logger        shared.ILogger

	appCtx    context.Context
	appCancel context.CancelFunc
}

func NewApplicationBuilder(configuration configuration.Configuration, startup IStartup) IApplicationBuilder {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	return &Application{
		configuration: configuration,
		startup:       startup,

		appCtx:    ctx,
		appCancel: cancel,
	}
}

func (a *Application) SetServiceProvider(serviceProvider shared.IServiceProvider) IApplicationBuilder {
	a.container = serviceProvider
	return a
}

func (a *Application) SetWebHost(host shared.IHost) IApplicationBuilder {
	a.host = host
	return a
}

func (a *Application) Build() (IApplication, error) {
	var container shared.IServiceProvider
	var httpHandler http.Handler
	var logger shared.ILogger
	var err error

	if a.container == nil {
		container, err = ioc.NewDigServiceProvider()
		if err != nil {
			a.logger.Fatal("error occurred during IoC construction ")
		}
		a.container = container
	}

	err = a.startup.ConfigureServices(a.appCtx, a.configuration, a.container)
	if err != nil {
		a.logger.Fatal("error occurred during service provider configuration")
	}

	if a.host == nil {
		if err = a.container.GetService(func(h http.Handler) { httpHandler = h }); err != nil {
			a.logger.Fatal("error occurred during web server construction ")
		}
		srv := httpserver.NewServer(httpHandler, a.configuration)
		a.host = &srv
	}

	err = a.container.GetService(func(l shared.ILogger) { logger = l })

	if err != nil {
		return nil, err
	}

	a.logger = logger
	return a, err
}

func (a *Application) Run() {
	defer a.appCancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go a.handleShutdown(&wg)
	go func() {
		err := a.startup.AfterApplicationStartup(a.appCtx, a.container)
		if err != nil {
			a.logger.Fatal(fmt.Sprintf("after startup action execution ended with error %s\n", err.Error()))
		}
	}()

	a.logger.Info(fmt.Sprintf("Server started on %s:%d", a.configuration.Host, a.configuration.Port))
	if err := a.host.Start(a.appCtx); err != nil && err != http.ErrServerClosed {
		a.logger.Error(fmt.Sprintf("Error occurred while running http server: %s\n", err.Error()))
		a.appCancel()
		os.Exit(1) //nolint:gocritic
	}

	wg.Wait()
}

func (a *Application) handleShutdown(wg *sync.WaitGroup) {
	defer wg.Done()

	<-a.appCtx.Done()
	a.logger.Info("Start process of server shutdown ...")

	const timeout = 3 * time.Second

	var err error
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	defer func() {
		err = a.logger.Flush()
	}()

	_ = a.startup.BeforeApplicationShutdown(ctx, a.container)

	if err = a.host.Stop(ctx); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to stop server: %v", err))
	}

	a.logger.Info("Server is stopped!")
}
