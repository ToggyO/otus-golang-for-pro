package memorystorage

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/models"
	domain "github.com/ToggyO/otus-golang-for-pro/hw12_13_14_15_calendar/internal/domain/repositories"
)

func TestInMemoryEventsRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryEventsRepository()

	seedData(t, ctx, repo)

	t.Run("get events list", func(t *testing.T) {
		filter := &models.EventsFilter{
			BaseFilter: models.BaseFilter{
				Page:     0,
				PageSize: 5,
			},

			StartDate: time.Date(2022, time.Month(9), 11, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.Month(9), 13, 0, 0, 0, 0, time.UTC),
		}

		events, err := repo.GetEventsList(ctx, filter)
		require.NoError(t, err)
		require.Len(t, events, 3)
	})

	t.Run("get by id", func(t *testing.T) {
		var id int64 = 2
		event, err := repo.GetEventById(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, event)
		require.Equal(t, id, event.ID)
	})

	t.Run("concurrent update", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		reg := regexp.MustCompile("[0-9]+$")

		updateFunc := func(t *testing.T, ctx context.Context, wg *sync.WaitGroup, index int, repository domain.IEventsRepository) {
			defer wg.Done()
			event, err := repository.UpdateEvent(ctx, 1, &models.EventInfo{
				Title:            fmt.Sprintf("event%d", index),
				StartDate:        time.Date(2022, time.Month(9), 12, 0, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2022, time.Month(9), 12, 2, 0, 0, 0, time.UTC),
				Description:      "Desc",
				OwnerId:          10,
				NotificationDate: time.Date(2022, time.Month(9), 11, 10, 0, 0, 0, time.UTC),
			})
			require.NoError(t, err)

			s := reg.FindAllString(event.Title, -1)
			num, err := strconv.Atoi(strings.Join(s, ""))
			require.Equal(t, index, num)
		}

		for i := 1; i <= 15; i++ {
			wg.Add(1)
			i := i
			go updateFunc(t, ctx, wg, i, repo)
		}

		wg.Wait()
	})

	t.Run("remove by id", func(t *testing.T) {
		event, err := repo.CreateEvent(ctx, getEvent("event4"))
		require.NoError(t, err)

		err = repo.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)

		event, err = repo.GetEventById(ctx, event.ID)
		require.Error(t, err)
	})
}

func seedData(t *testing.T, ctx context.Context, repository domain.IEventsRepository) {
	t.Helper()

	eventsInfos := []*models.EventInfo{
		getEvent("event1"),
		getEvent("event2"),
		getEvent("event3"),
	}

	var err error
	for _, e := range eventsInfos {
		_, err = repository.CreateEvent(ctx, e)
		require.Nil(t, err)
	}
}

func getEvent(title string) *models.EventInfo {
	return &models.EventInfo{
		Title:            title,
		StartDate:        time.Date(2022, time.Month(9), 12, 0, 0, 0, 0, time.UTC),
		EndDate:          time.Date(2022, time.Month(9), 12, 2, 0, 0, 0, time.UTC),
		Description:      "Desc",
		OwnerId:          10,
		NotificationDate: time.Date(2022, time.Month(9), 11, 10, 0, 0, 0, time.UTC),
	}
}
