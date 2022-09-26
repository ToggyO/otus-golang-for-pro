package api

import (
	"context"
	"net/http"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/api/http/controllers"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/api/http/middleware"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
	"go.uber.org/dig"
)

type appControllers struct {
	dig.In
	List []controllers.IRouteBinder `group:"controllers"`
}

// TODO: придумать, как убрать импорт dig.
func AddHTTPHandler(_ context.Context, serviceProvider shared.IServiceProvider) error {
	err := serviceProvider.AddService(&shared.ServiceDescriptor{Service: http.NewServeMux})
	if err != nil {
		return err
	}

	options := make([]dig.ProvideOption, 0, 1)
	options = append(options, dig.Group("controllers"))
	err = serviceProvider.AddService(&shared.ServiceDescriptor{Service: controllers.NewTestController, Options: options})
	if err != nil {
		return err
	}

	var ctrls appControllers
	err = serviceProvider.GetService(func(ac appControllers) {
		ctrls = ac
	})
	if err != nil {
		return err
	}

	for _, c := range ctrls.List {
		c.Bind()
	}

	bindHTTPHandler := func(mux *http.ServeMux, logger shared.ILogger) http.Handler {
		router := middleware.LoggingMiddleware(mux, logger)
		return router
	}

	err = serviceProvider.AddService(&shared.ServiceDescriptor{Service: bindHTTPHandler})

	return err
}
