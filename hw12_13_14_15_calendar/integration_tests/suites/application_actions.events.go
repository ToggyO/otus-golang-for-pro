package suites

import (
	"time"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/stretchr/testify/require"
)

func (ap *ApplicationActionsSuite) InitEventsTests() {
	eventsInfos := []*models.EventInfo{
		getEvent("event1"),
		getEvent("event2"),
		getEvent("event3"),
	}

	var err error
	for _, e := range eventsInfos {
		_, err = ap.repositories.eventsRepository.CreateEvent(ap.ctx, e)
		require.NoError(ap.Suite.T(), err)
	}
}

func (ap *ApplicationActionsSuite) TestGetEventsList() {
	filter := &models.EventsFilter{
		BaseFilter: models.BaseFilter{
			Page:     0,
			PageSize: 5,
		},
	}

	events, err := ap.repositories.eventsRepository.GetEventsList(ap.ctx, filter)
	require.NoError(ap.Suite.T(), err)
	require.Len(ap.Suite.T(), events, 3)
}

func (ap *ApplicationActionsSuite) TestUpdateEvent() {
	var eventID int64 = 2
	eventTitle := "event69"
	newEventDesc := "Desc 69"

	eventInfo := getEvent(eventTitle)
	eventInfo.Description = newEventDesc

	_, err := ap.repositories.eventsRepository.UpdateEvent(ap.ctx, eventID, eventInfo)
	require.NoError(ap.Suite.T(), err)

	event, err := ap.repositories.eventsRepository.GetEventByID(ap.ctx, eventID)
	require.NoError(ap.Suite.T(), err)
	require.NotNil(ap.Suite.T(), event)
	require.Equal(ap.Suite.T(), eventID, event.ID)
	require.Equal(ap.Suite.T(), eventTitle, event.Title)
	require.Equal(ap.Suite.T(), newEventDesc, event.Description)
}

func (ap *ApplicationActionsSuite) TestGetEventByID() {
	var eventID int64 = 1

	event, err := ap.repositories.eventsRepository.GetEventByID(ap.ctx, eventID)
	require.NoError(ap.Suite.T(), err)
	require.NotNil(ap.Suite.T(), event)
	require.Equal(ap.Suite.T(), eventID, event.ID)
}

func (ap *ApplicationActionsSuite) TestDeleteEvent() {
	eventInfo := getEvent("event16")

	event, err := ap.repositories.eventsRepository.CreateEvent(ap.ctx, eventInfo)
	require.NoError(ap.Suite.T(), err)
	require.NotNil(ap.Suite.T(), event)

	err = ap.repositories.eventsRepository.DeleteEvent(ap.ctx, event.ID)
	require.NoError(ap.Suite.T(), err)

	event, err = ap.repositories.eventsRepository.GetEventByID(ap.ctx, event.ID)
	require.Error(ap.Suite.T(), err)
	require.Nil(ap.Suite.T(), event)
}

func getEvent(title string) *models.EventInfo {
	return &models.EventInfo{
		Title:            title,
		StartDate:        time.Date(2022, time.Month(9), 12, 0, 0, 0, 0, time.UTC),
		EndDate:          time.Date(2022, time.Month(9), 12, 2, 0, 0, 0, time.UTC),
		Description:      "Desc",
		OwnerID:          10,
		NotificationDate: time.Date(2022, time.Month(9), 11, 10, 0, 0, 0, time.UTC),
	}
}
