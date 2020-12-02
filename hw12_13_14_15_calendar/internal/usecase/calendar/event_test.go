package calendar

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/mocks"
	"github.com/stretchr/testify/require"
)

func TestEventUseCase(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		e := Event{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		}
		ctx := context.Background()
		storEvent := FromEvent(e)

		rep.On("CreateEvent", ctx, storEvent).
			Return(storage.EventID(1), nil)

		useCase := NewEventUseCase(rep)
		actualEvent, err := useCase.CreateEvent(ctx, e)

		require.NoError(t, err)
		require.Equal(t, e, actualEvent)
	})

	t.Run("create event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		var storEvent storage.Event

		rep.On("CreateEvent", ctx, storEvent).
			Return(storage.EventID(0), fmt.Errorf("create error"))

		useCase := NewEventUseCase(rep)
		_, err := useCase.CreateEvent(ctx, Event{})

		require.Error(t, err)
	})

	t.Run("get event by id", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		expected := storage.Event{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		}

		rep.On("GetEventByID", ctx, storage.EventID(1)).
			Return(expected, nil)

		useCase := NewEventUseCase(rep)
		actual, err := useCase.GetEventByID(ctx, 1)

		require.NoError(t, err)
		require.Equal(t, ToEvent(expected), actual)
	})

	t.Run("get event by id error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		rep.On("GetEventByID", ctx, storage.EventID(1)).
			Return(storage.Event{}, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		_, err := useCase.GetEventByID(ctx, 1)

		require.Error(t, err)
	})

	t.Run("update event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		rep.On("UpdateEvent", ctx, storage.Event{ID: 1}).
			Return(nil)

		useCase := NewEventUseCase(rep)
		err := useCase.UpdateEvent(ctx, 1, Event{})

		require.NoError(t, err)
	})

	t.Run("update event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		rep.On("UpdateEvent", ctx, storage.Event{ID: 1}).
			Return(fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		err := useCase.UpdateEvent(ctx, 1, Event{})

		require.Error(t, err)
	})

	t.Run("delete event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		rep.On("DeleteEvent", ctx, storage.EventID(1)).
			Return(nil)

		useCase := NewEventUseCase(rep)
		err := useCase.DeleteEvent(ctx, 1)

		require.NoError(t, err)
	})

	t.Run("delete event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		rep.On("DeleteEvent", ctx, storage.EventID(1)).
			Return(fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		err := useCase.DeleteEvent(ctx, 1)

		require.Error(t, err)
	})

	t.Run("get user day events", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		storEvents := []storage.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		sDate := now.With(curTime).BeginningOfDay()
		eDate := now.With(curTime).EndOfDay()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(storEvents, nil)

		useCase := NewEventUseCase(rep)
		actualEvents, err := useCase.GetUserDayEvents(ctx, 1, curTime)

		require.NoError(t, err)
		require.Equal(t, ToEventSlice(storEvents), actualEvents)
	})

	t.Run("get user day events error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		sDate := now.With(curTime).BeginningOfDay()
		eDate := now.With(curTime).EndOfDay()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(nil, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		_, err := useCase.GetUserDayEvents(ctx, 1, curTime)

		require.Error(t, err)
	})

	t.Run("get user week events", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		storEvents := []storage.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		sDate := now.With(curTime).BeginningOfWeek()
		eDate := now.With(curTime).EndOfWeek()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(storEvents, nil)

		useCase := NewEventUseCase(rep)
		actualEvents, err := useCase.GetUserWeekEvents(ctx, 1, curTime)

		require.NoError(t, err)
		require.Equal(t, ToEventSlice(storEvents), actualEvents)
	})

	t.Run("get user week events error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		sDate := now.With(curTime).BeginningOfWeek()
		eDate := now.With(curTime).EndOfWeek()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(nil, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		_, err := useCase.GetUserWeekEvents(ctx, 1, curTime)

		require.Error(t, err)
	})

	t.Run("get user month events", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		storEvents := []storage.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		sDate := now.With(curTime).BeginningOfMonth()
		eDate := now.With(curTime).EndOfMonth()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(storEvents, nil)

		useCase := NewEventUseCase(rep)
		actualEvents, err := useCase.GetUserMonthEvents(ctx, 1, curTime)

		require.NoError(t, err)
		require.Equal(t, ToEventSlice(storEvents), actualEvents)
	})

	t.Run("get user month events error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		curTime := time.Now()

		sDate := now.With(curTime).BeginningOfMonth()
		eDate := now.With(curTime).EndOfMonth()

		ctx := context.Background()
		rep.On("GetUserEventsByPeriod", ctx, storage.UserID(1), sDate, eDate).
			Return(nil, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		_, err := useCase.GetUserMonthEvents(ctx, 1, curTime)

		require.Error(t, err)
	})
}

func TestToEvent(t *testing.T) {
	se := storage.Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        time.Now(),
		EndDate:          time.Now(),
		NotificationDate: time.Now(),
	}

	expected := Event{
		ID:               int64(se.ID),
		Title:            se.Title,
		Description:      se.Description,
		UserID:           int64(se.UserID),
		StartDate:        se.StartDate,
		EndDate:          se.EndDate,
		NotificationDate: se.NotificationDate,
	}

	require.Equal(t, expected, ToEvent(se))
}

func TestFromEvent(t *testing.T) {
	e := Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        time.Now(),
		EndDate:          time.Now(),
		NotificationDate: time.Now(),
	}

	expected := storage.Event{
		ID:               storage.EventID(e.ID),
		Title:            e.Title,
		Description:      e.Description,
		UserID:           storage.UserID(e.UserID),
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		NotificationDate: e.NotificationDate,
	}

	require.Equal(t, expected, FromEvent(e))
}

func TestToEventSlice(t *testing.T) {
	se := []storage.Event{
		{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		},
		{
			ID:               2,
			Title:            "title2",
			Description:      "description2",
			UserID:           2,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		},
	}

	expected := []Event{
		{
			ID:               int64(se[0].ID),
			Title:            se[0].Title,
			Description:      se[0].Description,
			UserID:           int64(se[0].UserID),
			StartDate:        se[0].StartDate,
			EndDate:          se[0].EndDate,
			NotificationDate: se[0].NotificationDate,
		},
		{
			ID:               int64(se[1].ID),
			Title:            se[1].Title,
			Description:      se[1].Description,
			UserID:           int64(se[1].UserID),
			StartDate:        se[1].StartDate,
			EndDate:          se[1].EndDate,
			NotificationDate: se[1].NotificationDate,
		},
	}

	require.Equal(t, expected, ToEventSlice(se))
}
