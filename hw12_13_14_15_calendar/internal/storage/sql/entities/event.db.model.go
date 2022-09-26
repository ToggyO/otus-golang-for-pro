package entities

import (
	"database/sql"
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/apperrors"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
)

type EventDBModel struct {
	ID int64 `db:"id"`

	BaseDBModel       `db:""`
	*EventInfoDBModel `db:""`
}

type EventInfoDBModel struct {
	Title            string         `db:"title"`
	StartDate        time.Time      `db:"start_date"`
	EndDate          sql.NullTime   `db:"end_date"`
	Description      sql.NullString `db:"description"`
	OwnerID          int64          `db:"owner_id"`
	NotificationDate time.Time      `db:"notification_date"`
}

func (e *EventInfoDBModel) FromEventInfo(eventInfo *models.EventInfo) error {
	if eventInfo == nil {
		return apperrors.ErrArgumentNil
	}

	e.Title = eventInfo.Title
	e.StartDate = eventInfo.StartDate

	e.EndDate = sql.NullTime{
		Time:  eventInfo.EndDate,
		Valid: !eventInfo.EndDate.IsZero(),
	}

	e.Description = sql.NullString{
		String: eventInfo.Description,
		Valid:  eventInfo.Description != "",
	}

	e.OwnerID = eventInfo.OwnerID
	e.NotificationDate = eventInfo.NotificationDate

	return nil
}
