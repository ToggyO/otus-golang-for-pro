package app

import (
	appcore "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/app/core"
	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
)

type serviceProvider struct {
	db               appcore.IStorage
	eventsRepository domain.IEventsRepository
}
