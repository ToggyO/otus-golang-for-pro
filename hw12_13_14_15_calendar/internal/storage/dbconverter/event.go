package dbconverter

import (
	"database/sql"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/apperrors"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/storage/sql/entities"
)

func ToEventInfo(dbEventInfo *entities.EventInfoDbModel) (*models.EventInfo, error) {
	if dbEventInfo == nil {
		return nil, apperrors.ErrArgumentNil
	}

	eventInfo := &models.EventInfo{
		Title:            dbEventInfo.Title,
		StartDate:        dbEventInfo.StartDate,
		EndDate:          dbEventInfo.EndDate.Time,
		Description:      dbEventInfo.Description.String,
		OwnerId:          dbEventInfo.OwnerId,
		NotificationDate: dbEventInfo.NotificationDate,
	}

	return eventInfo, nil
}

func FromEventInfo(eventInfo *models.EventInfo) (*entities.EventInfoDbModel, error) {
	if eventInfo == nil {
		return nil, apperrors.ErrArgumentNil
	}

	dbEventInfo := &entities.EventInfoDbModel{
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

		OwnerId:          eventInfo.OwnerId,
		NotificationDate: eventInfo.NotificationDate,
	}

	return dbEventInfo, nil
}

func ToEvent(eventDbModel *entities.EventDbModel) (*models.Event, error) {
	if eventDbModel == nil {
		return nil, apperrors.ErrArgumentNil
	}

	eventInfo, err := ToEventInfo(eventDbModel.EventInfoDbModel)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		ID:        eventDbModel.ID,
		EventInfo: eventInfo,
	}

	return event, nil
}
