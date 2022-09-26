package entities

import (
	"database/sql"
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/apperrors"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
)

type EventDbModel struct {
	ID int64 `db:"id"`

	BaseDbModel       `db:""`
	*EventInfoDbModel `db:""`
}

type EventInfoDbModel struct {
	Title            string         `db:"title"`
	StartDate        time.Time      `db:"start_date"`
	EndDate          sql.NullTime   `db:"end_date"`
	Description      sql.NullString `db:"description"`
	OwnerId          int64          `db:"owner_id"`
	NotificationDate time.Time      `db:"notification_date"`
}

func (e EventInfoDbModel) FromEventInfo(eventInfo *models.EventInfo) error {
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

	e.OwnerId = eventInfo.OwnerId
	e.NotificationDate = eventInfo.NotificationDate

	return nil
}
