package ioc

import (
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/dig"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/pkg/shared"
)

type serviceProvider struct {
	container *dig.Container
}

func NewDigServiceProvider() (shared.IServiceProvider, error) {
	return &serviceProvider{container: dig.New()}, nil
}

func (c *serviceProvider) GetService(function interface{}) error {
	err := c.container.Invoke(function)
	return err
}

func (c *serviceProvider) AddService(sd *shared.ServiceDescriptor) error {
	if sd.Options == nil {
		return c.container.Provide(sd.Service)
	}
	return c.container.Provide(sd.Service, sd.Options.([]dig.ProvideOption)...)
}

func (c *serviceProvider) RunAfterBuild(functionList []func()) error {
	for _, f := range functionList {
		fType := reflect.TypeOf(f)
		if fType == nil {
			return errors.New("Can't invoke an untyped nil ")
		}

		if fType.Kind() != reflect.Func {
			return errors.New(fmt.Sprintf("can't invoke non-function (type %v)", fType))
		}

		f()
	}

	return nil
}
