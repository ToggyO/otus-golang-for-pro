package dbconverter

import (
	"database/sql"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/apperrors"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/entities"
)

func ToEventInfo(dbEventInfo *entities.EventInfoDBModel) (*models.EventInfo, error) {
	if dbEventInfo == nil {
		return nil, apperrors.ErrArgumentNil
	}

	eventInfo := &models.EventInfo{
		Title:            dbEventInfo.Title,
		StartDate:        dbEventInfo.StartDate,
		EndDate:          dbEventInfo.EndDate.Time,
		Description:      dbEventInfo.Description.String,
		OwnerID:          dbEventInfo.OwnerID,
		NotificationDate: dbEventInfo.NotificationDate,
	}

	return eventInfo, nil
}

func FromEventInfo(eventInfo *models.EventInfo) (*entities.EventInfoDBModel, error) {
	if eventInfo == nil {
		return nil, apperrors.ErrArgumentNil
	}

	dbEventInfo := &entities.EventInfoDBModel{
		Title:     eventInfo.Title,
		StartDate: eventInfo.StartDate,

		EndDate: sql.NullTime{
			Time:  eventInfo.EndDate,
			Valid: !eventInfo.EndDate.IsZero(),
		},

		Description: sql.NullString{
			String: eventInfo.Description,
			Valid:  eventInfo.Description != "",
		},

		OwnerID:          eventInfo.OwnerID,
		NotificationDate: eventInfo.NotificationDate,
	}

	return dbEventInfo, nil
}

func ToEvent(eventDBModel *entities.EventDBModel) (*models.Event, error) {
	if eventDBModel == nil {
		return nil, apperrors.ErrArgumentNil
	}

	eventInfo, err := ToEventInfo(eventDBModel.EventInfoDBModel)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		ID:        eventDBModel.ID,
		EventInfo: eventInfo,
	}

	return event, nil
}
