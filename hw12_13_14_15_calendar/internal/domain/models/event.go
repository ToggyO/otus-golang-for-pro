package models

import (
	"time"
)

type Event struct {
	ID int64

	*EventInfo
}

func (e Event) Clone() *Event {
	return &Event{
		ID:        e.ID,
		EventInfo: e.EventInfo.Clone(),
	}
}

type EventInfo struct {
	Title            string
	StartDate        time.Time
	EndDate          time.Time
	Description      string
	OwnerId          int64
	NotificationDate time.Time
}

func (ei EventInfo) Clone() *EventInfo {
	return &EventInfo{
		Title:            ei.Title,
		StartDate:        ei.StartDate,
		EndDate:          ei.EndDate,
		Description:      ei.Description,
		OwnerId:          ei.OwnerId,
		NotificationDate: ei.NotificationDate,
	}
}

type EventsFilter struct {
	BaseFilter

	StartDate time.Time
	EndDate   time.Time
}
