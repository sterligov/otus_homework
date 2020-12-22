package calendar

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/mocks"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/model"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestEventUseCase_CreateEvent(t *testing.T) {
	t.Run("create event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		e := model.Event{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		}
		ctx := context.Background()
		storEvent := model.FromEvent(e)

		rep.On("CreateEvent", ctx, storEvent).
			Return(storage.EventID(1), nil)

		useCase := NewEventUseCase(rep)
		insertedID, err := useCase.CreateEvent(ctx, e)

		require.NoError(t, err)
		require.Equal(t, int64(1), insertedID)
	})

	t.Run("create event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		ctx := context.Background()
		var storEvent storage.Event

		rep.On("CreateEvent", ctx, storEvent).
			Return(storage.EventID(0), fmt.Errorf("create error"))

		useCase := NewEventUseCase(rep)
		insertedID, err := useCase.CreateEvent(ctx, model.Event{})

		require.Error(t, err)
		require.Equal(t, int64(0), insertedID)
	})
}

func TestEventUseCase_GetEventByID(t *testing.T) {
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
		require.Equal(t, model.ToEvent(expected), actual)
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
}

func TestEventUseCase_DeleteEvent(t *testing.T) {
	t.Run("delete event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		expectedAffected := int64(1)

		ctx := context.Background()
		rep.On("DeleteEvent", ctx, storage.EventID(1)).
			Return(expectedAffected, nil)

		useCase := NewEventUseCase(rep)
		affected, err := useCase.DeleteEvent(ctx, 1)

		require.NoError(t, err)
		require.Equal(t, expectedAffected, affected)
	})

	t.Run("delete event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		var expectedAffected int64

		ctx := context.Background()
		rep.On("DeleteEvent", ctx, storage.EventID(1)).
			Return(expectedAffected, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		affected, err := useCase.DeleteEvent(ctx, 1)

		require.Error(t, err)
		require.Equal(t, expectedAffected, affected)
	})
}

func TestEventUseCase_UpdateEvent(t *testing.T) {
	t.Run("update event", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		expectedAffected := int64(1)

		ctx := context.Background()
		rep.On("UpdateEvent", ctx, storage.Event{ID: 1}).
			Return(expectedAffected, nil)

		useCase := NewEventUseCase(rep)
		affected, err := useCase.UpdateEvent(ctx, 1, model.Event{})

		require.NoError(t, err)
		require.Equal(t, expectedAffected, affected)
	})

	t.Run("update event error", func(t *testing.T) {
		rep := &mocks.EventRepository{}

		var expectedAffected int64

		ctx := context.Background()
		rep.On("UpdateEvent", ctx, storage.Event{ID: 1}).
			Return(expectedAffected, fmt.Errorf("error here"))

		useCase := NewEventUseCase(rep)
		affected, err := useCase.UpdateEvent(ctx, 1, model.Event{})

		require.Error(t, err)
		require.Equal(t, expectedAffected, affected)
	})
}

func TestEventUseCase_GetUserDayEvents(t *testing.T) {
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
		require.Equal(t, model.ToEventSlice(storEvents), actualEvents)
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
}

func TestEventUseCase_GetUserWeekEvents(t *testing.T) {
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
		require.Equal(t, model.ToEventSlice(storEvents), actualEvents)
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
}

func TestEventUseCase_GetUserMonthEvents(t *testing.T) {
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
		require.Equal(t, model.ToEventSlice(storEvents), actualEvents)
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
