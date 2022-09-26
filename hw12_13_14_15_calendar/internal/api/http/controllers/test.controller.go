package controllers

import (
	"net/http"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type TestController struct {
	router *http.ServeMux
	logger shared.ILogger
}

func NewTestController(router *http.ServeMux, logger shared.ILogger) IRouteBinder {
	return &TestController{
		router: router,
		logger: logger,
	}
}

func (tc *TestController) Bind() {
	tc.router.HandleFunc(routesTestGet, tc.Ping)
}

func (tc *TestController) Ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Pong"))
}
