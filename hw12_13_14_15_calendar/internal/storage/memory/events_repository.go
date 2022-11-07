package memorystorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/apperrors"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
)

type eventsRepository struct {
	mu             *sync.RWMutex
	events         map[int64]*models.Event
	lastInsertedID int64
}

func NewInMemoryEventsRepository() domain.IEventsRepository {
	return &eventsRepository{
		mu:     &sync.RWMutex{},
		events: make(map[int64]*models.Event),
	}
}

func (e *eventsRepository) CreateEvent(_ context.Context, eventInfo *models.EventInfo) (*models.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.lastInsertedID++
	event := &models.Event{
		ID:        e.lastInsertedID,
		EventInfo: eventInfo,
	}

	e.events[e.lastInsertedID] = event

	// Prevent mutation of original event in memory storage from outer scope
	return event.Clone(), nil
}

func (e *eventsRepository) UpdateEvent(
	_ context.Context,
	id int64,
	eventInfo *models.EventInfo,
) (*models.Event, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	event, ok := e.events[id]
	if !ok {
		return nil, apperrors.ErrNotFound(fmt.Sprintf("event with id=%d", id))
	}

	event.EventInfo = eventInfo
	e.events[id] = event

	// Prevent mutation of original event in memory storage from outer scope
	return event.Clone(), nil
}

func (e *eventsRepository) DeleteEvent(_ context.Context, id int64) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.events, id)
	return nil
}

func (e *eventsRepository) GetEventsList(_ context.Context, filter *models.EventsFilter) ([]models.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	eventsList := make([]models.Event, 0, filter.PageSize)

	for index, ev := range e.events {
		if uint64(len(eventsList)) > filter.PageSize {
			break
		}

		if uint64(index) < filter.Page {
			continue
		}

		if ev.StartDate.After(filter.StartDate) && ev.EndDate.Before(filter.EndDate) {
			eventsList = append(eventsList, *ev)
		}
	}

	return eventsList, nil
}

func (e *eventsRepository) GetEventByID(_ context.Context, id int64) (*models.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	event, ok := e.events[id]
	if !ok {
		return nil, apperrors.ErrNotFound(fmt.Sprintf("event with id=%d", id))
	}

	// Prevent mutation of original event in memory storage from outer scope
	return event.Clone(), nil
}
