package memorystorage

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

const dateLayout = "2006-01-02 15:04"

func TestEventStorage(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		stor := NewEventStorage()
		ctx := context.Background()

		e := storage.Event{}

		insertedID, err := stor.CreateEvent(ctx, e)
		require.NoError(t, err)
		require.EqualValues(t, 1, insertedID)
	})

	t.Run("get", func(t *testing.T) {
		stor := NewEventStorage()
		ctx := context.Background()

		expected := storage.Event{
			UserID:           1,
			Title:            "title",
			Description:      "description",
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		}

		insertedID, err := stor.CreateEvent(ctx, expected)
		require.NoError(t, err)

		expected.ID = insertedID
		actual, err := stor.GetEventByID(ctx, insertedID)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("delete", func(t *testing.T) {
		stor := NewEventStorage()
		ctx := context.Background()

		e := storage.Event{ID: 1}

		_, err := stor.CreateEvent(ctx, e)
		require.NoError(t, err)

		affected, err := stor.DeleteEvent(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)

		affected, err = stor.DeleteEvent(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, int64(0), affected)
	})

	t.Run("update", func(t *testing.T) {
		stor := NewEventStorage()

		e := storage.Event{ID: 1}

		_, err := stor.CreateEvent(context.Background(), e)
		require.NoError(t, err)

		affected, err := stor.UpdateEvent(context.Background(), e)
		require.NoError(t, err)
		require.Equal(t, int64(1), affected)
	})

	t.Run("delete before date", func(t *testing.T) {
		stor := NewEventStorage()
		ctx := context.Background()

		e := storage.Event{ID: 1, StartDate: time.Now().Add(-time.Minute)}

		insertedID, err := stor.CreateEvent(ctx, e)
		require.NoError(t, err)

		_, err = stor.DeleteEventsBeforeDate(ctx, time.Now().Add(-2*time.Minute))
		require.NoError(t, err)
		_, err = stor.GetEventByID(ctx, insertedID)
		require.NoError(t, nil)

		_, err = stor.DeleteEventsBeforeDate(ctx, time.Now())
		require.NoError(t, err)
		_, err = stor.GetEventByID(ctx, insertedID)
		require.True(t, errors.Is(err, storage.ErrNotFound))
	})

	t.Run("get by notification date period", func(t *testing.T) {
		stor := NewEventStorage()
		ctx := context.Background()

		e := storage.Event{ID: 1, StartDate: time.Now().Add(-time.Minute)}

		_, err := stor.CreateEvent(ctx, e)
		require.NoError(t, err)

		events, err := stor.GetEventsByNotificationDatePeriod(ctx, time.Now().Add(-3*time.Minute), time.Now().Add(-2*time.Minute))
		require.NoError(t, err)
		require.Empty(t, events)

		events, err = stor.GetEventsByNotificationDatePeriod(ctx, time.Now().Add(-3*time.Minute), time.Now())
		require.NoError(t, err)
		require.NotEmpty(t, events)
	})

	t.Run("create two events in one date", func(t *testing.T) {
		stor := NewEventStorage()

		curTime := time.Now()
		e1 := storage.Event{UserID: 1, StartDate: curTime}
		e2 := storage.Event{UserID: 1, StartDate: curTime}
		e3 := storage.Event{UserID: 2, StartDate: curTime}

		_, err := stor.CreateEvent(context.Background(), e1)
		require.NoError(t, err)

		_, err = stor.CreateEvent(context.Background(), e2)
		require.Equal(t, storage.ErrDateBusy, err)

		_, err = stor.CreateEvent(context.Background(), e3)
		require.NoError(t, err)
	})

	t.Run("update event and set existing date", func(t *testing.T) {
		stor := NewEventStorage()

		curTime := time.Now()
		e1 := storage.Event{UserID: 1, StartDate: curTime.Round(0)}
		e2 := storage.Event{UserID: 1, StartDate: curTime.Add(time.Hour).Round(0)}

		insertedID, err := stor.CreateEvent(context.Background(), e1)
		require.NoError(t, err)
		e1.ID = insertedID

		_, err = stor.CreateEvent(context.Background(), e2)
		require.NoError(t, err)

		e1.StartDate = e1.StartDate.Add(time.Hour).Round(0)
		_, err = stor.UpdateEvent(context.Background(), e1)
		require.Equal(t, storage.ErrDateBusy, err)
	})

	t.Run("not found", func(t *testing.T) {
		stor := NewEventStorage()

		_, err := stor.GetEventByID(context.Background(), 1)
		require.Equal(t, storage.ErrNotFound, err)
	})

	t.Run("complex", func(t *testing.T) {
		events := []storage.Event{
			{
				Title:            "title",
				Description:      "description",
				UserID:           100,
				StartDate:        string2Time(t, "2020-12-01 10:00"),
				EndDate:          string2Time(t, "2020-12-01 10:05"),
				NotificationDate: string2Time(t, "2020-12-02 09:55"),
			},
			{
				Title:            "title2",
				Description:      "description2",
				UserID:           100,
				StartDate:        string2Time(t, "2020-12-01 17:00"),
				EndDate:          string2Time(t, "2020-12-01 17:15"),
				NotificationDate: string2Time(t, "2020-12-01 16:50"),
			},
			{
				Title:            "title3",
				Description:      "description3",
				UserID:           100,
				StartDate:        string2Time(t, "2020-12-02 15:00"),
				EndDate:          string2Time(t, "2020-12-02 16:05"),
				NotificationDate: string2Time(t, "2020-12-02 14:55"),
			},
			{
				Title:            "title4",
				Description:      "description4",
				UserID:           200,
				StartDate:        string2Time(t, "2020-12-10 18:00"),
				EndDate:          string2Time(t, "2020-12-10 18:10"),
				NotificationDate: string2Time(t, "2020-12-10 17:55"),
			},
		}

		stor := NewEventStorage()
		ctx := context.Background()
		var wg sync.WaitGroup

		nEvent := len(events)
		wg.Add(nEvent)
		for i := 0; i < nEvent; i++ {
			go func(i int) {
				defer wg.Done()

				insertedID, err := stor.CreateEvent(ctx, events[i])
				require.NoError(t, err)
				events[i].ID = insertedID
			}(i)
		}
		wg.Wait()

		start := string2Time(t, "2020-12-01 00:00")
		end := string2Time(t, "2020-12-31 00:00")

		actualEvents, err := stor.GetUserEventsByPeriod(ctx, events[0].UserID, start, end)
		require.NoError(t, err)
		require.ElementsMatch(t, events[0:3], actualEvents)

		wg.Add(2)
		go func() {
			defer wg.Done()

			events[0].Title = "New Title"
			_, err := stor.UpdateEvent(ctx, events[0])
			require.NoError(t, err)
		}()
		go func() {
			defer wg.Done()

			_, err := stor.DeleteEvent(ctx, events[1].ID)
			require.NoError(t, err)
		}()
		wg.Wait()

		start = string2Time(t, "2020-12-01 00:00")
		end = string2Time(t, "2020-12-02 00:00")
		actualEvents, err = stor.GetUserEventsByPeriod(ctx, events[0].UserID, start, end)
		require.NoError(t, err)
		require.ElementsMatch(t, events[0:1], actualEvents)
	})
}

func string2Time(t *testing.T, date string) time.Time {
	d, err := time.Parse(dateLayout, date)
	require.NoError(t, err)

	return d
}
