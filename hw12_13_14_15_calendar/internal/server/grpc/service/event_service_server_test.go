package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/mocks"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/model"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEventServiceServer_CreateEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		e := &pb.Event{
			Id:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        timestamppb.Now(),
			EndDate:          timestamppb.Now(),
			NotificationDate: timestamppb.Now(),
		}
		ctx := context.Background()
		insertedID := int64(1)

		eventUseCase.On("CreateEvent", ctx, FromEvent(e)).
			Return(insertedID, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		inserted, err := server.CreateEvent(context.Background(), e)

		require.NoError(t, err)
		require.Equal(t, insertedID, inserted.InsertedID)
	})

	t.Run("busy date", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		e := &pb.Event{}

		eventUseCase.On("CreateEvent", ctx, FromEvent(e)).
			Return(int64(0), storage.ErrDateBusy)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		inserted, err := server.CreateEvent(ctx, e)
		s, ok := status.FromError(err)

		require.Nil(t, inserted)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, s.Code())
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		e := &pb.Event{}

		eventUseCase.On("CreateEvent", ctx, FromEvent(e)).
			Return(int64(0), fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		inserted, err := server.CreateEvent(ctx, e)

		require.Error(t, err)
		require.Nil(t, inserted)
	})
}

func TestEventServiceServer_DeleteEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		affectedRows := int64(1)
		eventID := int64(1)

		eventUseCase.On("DeleteEvent", ctx, eventID).
			Return(affectedRows, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		affected, err := server.DeleteEvent(ctx, &pb.EventID{Id: eventID})

		require.NoError(t, err)
		require.Equal(t, affectedRows, affected.Affected)
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		eventID := int64(1)

		eventUseCase.On("DeleteEvent", ctx, eventID).
			Return(int64(0), fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		affected, err := server.DeleteEvent(ctx, &pb.EventID{Id: eventID})

		require.Error(t, err)
		require.Nil(t, affected)
	})
}

func TestEventServiceServer_UpdateEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		e := &pb.Event{}
		affectedRows := int64(1)

		ctx := context.Background()
		eventUseCase.On("UpdateEvent", ctx, e.Id, FromEvent(e)).
			Return(affectedRows, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		affected, err := server.UpdateEvent(ctx, e)

		require.NoError(t, err)
		require.Equal(t, affectedRows, affected.Affected)
	})

	t.Run("busy date", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		e := &pb.Event{}

		eventUseCase.On("UpdateEvent", ctx, e.Id, FromEvent(e)).
			Return(int64(0), storage.ErrDateBusy)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		inserted, err := server.UpdateEvent(ctx, e)
		s, ok := status.FromError(err)

		require.Nil(t, inserted)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, s.Code())
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		ctx := context.Background()
		e := &pb.Event{}

		eventUseCase.On("UpdateEvent", ctx, e.Id, FromEvent(e)).
			Return(int64(0), fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		inserted, err := server.UpdateEvent(ctx, e)

		require.Error(t, err)
		require.Nil(t, inserted)
	})
}

func TestEventServiceServer_GetEventByID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		eventID := int64(1)
		e := &pb.Event{
			Id:               eventID,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        timestamppb.Now(),
			EndDate:          timestamppb.Now(),
			NotificationDate: timestamppb.Now(),
		}

		ctx := context.Background()
		eventUseCase.On("GetEventByID", ctx, eventID).
			Return(FromEvent(e), nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvent, err := server.GetEventByID(ctx, &pb.EventID{Id: eventID})

		require.NoError(t, err)
		require.Equal(t, e, actualEvent)
	})

	t.Run("not found", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		eventID := int64(1)

		ctx := context.Background()
		eventUseCase.On("GetEventByID", ctx, eventID).
			Return(model.Event{}, storage.ErrNotFound)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvent, err := server.GetEventByID(ctx, &pb.EventID{Id: eventID})
		s, ok := status.FromError(err)

		require.Nil(t, actualEvent)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, s.Code())
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		eventID := int64(1)

		ctx := context.Background()
		eventUseCase.On("GetEventByID", ctx, eventID).
			Return(model.Event{}, fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvent, err := server.GetEventByID(ctx, &pb.EventID{Id: eventID})

		require.Nil(t, actualEvent)
		require.Error(t, err)
	})
}

func TestEventServiceServer_GetUserDayEvents(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)
		events := []model.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		ctx := context.Background()
		eventUseCase.On("GetUserDayEvents", ctx, userID, curTime.AsTime()).
			Return(events, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserDayEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.NoError(t, err)
		require.Equal(t, ToEvents(events), actualEvents)
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)

		ctx := context.Background()
		eventUseCase.On("GetUserDayEvents", ctx, userID, curTime.AsTime()).
			Return(nil, fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserDayEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.Error(t, err)
		require.Nil(t, actualEvents)
	})
}

func TestEventServiceServer_GetUserWeekEvents(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)
		events := []model.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		ctx := context.Background()
		eventUseCase.On("GetUserWeekEvents", ctx, userID, curTime.AsTime()).
			Return(events, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserWeekEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.NoError(t, err)
		require.Equal(t, ToEvents(events), actualEvents)
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)

		ctx := context.Background()
		eventUseCase.On("GetUserWeekEvents", ctx, userID, curTime.AsTime()).
			Return(nil, fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserWeekEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.Error(t, err)
		require.Nil(t, actualEvents)
	})
}

func TestEventServiceServer_GetUserMonthEvents(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)
		events := []model.Event{
			{
				ID:    1,
				Title: "title1",
			},
			{
				ID:    2,
				Title: "title2",
			},
		}

		ctx := context.Background()
		eventUseCase.On("GetUserMonthEvents", ctx, userID, curTime.AsTime()).
			Return(events, nil)

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserMonthEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.NoError(t, err)
		require.Equal(t, ToEvents(events), actualEvents)
	})

	t.Run("error", func(t *testing.T) {
		eventUseCase := &mocks.EventUseCase{}
		curTime := timestamppb.Now()
		userID := int64(1)

		ctx := context.Background()
		eventUseCase.On("GetUserMonthEvents", ctx, userID, curTime.AsTime()).
			Return(nil, fmt.Errorf("internal error"))

		server := NewEventServiceServer(eventUseCase, &mocks.StorageConnection{})
		actualEvents, err := server.GetUserMonthEvents(ctx, &pb.UserPeriodEventRequest{UserID: userID, Date: curTime})

		require.Error(t, err)
		require.Nil(t, actualEvents)
	})
}

func TestEventServiceServer_Health(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		storageConnection := &mocks.StorageConnection{}

		storageConnection.On("PingContext", mock.Anything).
			Return(nil)

		server := NewEventServiceServer(&mocks.EventUseCase{}, storageConnection)
		resp, err := server.Health(context.Background(), &pb.Empty{})

		require.NoError(t, err)
		require.NotEmpty(t, resp)
	})

	t.Run("error", func(T *testing.T) {
		storageConnection := &mocks.StorageConnection{}

		storageConnection.On("PingContext", mock.Anything).
			Return(fmt.Errorf("internal error"))

		server := NewEventServiceServer(&mocks.EventUseCase{}, storageConnection)
		resp, err := server.Health(context.Background(), &pb.Empty{})

		require.True(t, errors.Is(err, ErrServiceUnavailable))
		require.Nil(t, resp)
	})
}

func TestToEvent(t *testing.T) {
	curTime := timestamppb.Now()
	e := model.Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        curTime.AsTime(),
		EndDate:          curTime.AsTime(),
		NotificationDate: curTime.AsTime(),
	}
	expected := &pb.Event{
		Id:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        curTime,
		EndDate:          curTime,
		NotificationDate: curTime,
	}
	actual := ToEvent(e)

	require.Equal(t, expected, actual)
}

func TestFromEvent(t *testing.T) {
	curTime := timestamppb.Now()
	e := &pb.Event{
		Id:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        curTime,
		EndDate:          curTime,
		NotificationDate: curTime,
	}
	expected := model.Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        curTime.AsTime(),
		EndDate:          curTime.AsTime(),
		NotificationDate: curTime.AsTime(),
	}
	actual := FromEvent(e)

	require.Equal(t, expected, actual)
}

func TestToEvents(t *testing.T) {
	curTime := timestamppb.Now()
	e := []model.Event{
		{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        curTime.AsTime(),
			EndDate:          curTime.AsTime(),
			NotificationDate: curTime.AsTime(),
		},
		{
			ID:               2,
			Title:            "title2",
			Description:      "description2",
			UserID:           2,
			StartDate:        curTime.AsTime(),
			EndDate:          curTime.AsTime(),
			NotificationDate: curTime.AsTime(),
		},
	}
	expected := &pb.Events{
		Events: []*pb.Event{
			{
				Id:               1,
				Title:            "title",
				Description:      "description",
				UserID:           1,
				StartDate:        curTime,
				EndDate:          curTime,
				NotificationDate: curTime,
			},
			{
				Id:               2,
				Title:            "title2",
				Description:      "description2",
				UserID:           2,
				StartDate:        curTime,
				EndDate:          curTime,
				NotificationDate: curTime,
			},
		},
	}

	actual := ToEvents(e)

	require.Equal(t, expected, actual)
}
