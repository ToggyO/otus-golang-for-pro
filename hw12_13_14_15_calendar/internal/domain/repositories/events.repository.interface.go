package domain

import (
	"context"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
)

type IEventsRepository interface {
	CreateEvent(ctx context.Context, eventInfo *models.EventInfo) (*models.Event, error)
	UpdateEvent(ctx context.Context, id int64, eventInfo *models.EventInfo) (*models.Event, error)
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsList(ctx context.Context, filter *models.EventsFilter) ([]models.Event, error)
	GetEventByID(ctx context.Context, id int64) (*models.Event, error)
}
