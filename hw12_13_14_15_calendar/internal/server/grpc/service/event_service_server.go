package service

import (
	"context"
	"errors"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/model"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrServiceUnavailable = status.Error(codes.Unavailable, "not alive")

type (
	EventUseCase interface {
		GetEventByID(ctx context.Context, id int64) (model.Event, error)
		CreateEvent(ctx context.Context, e model.Event) (int64, error)
		UpdateEvent(ctx context.Context, id int64, e model.Event) (int64, error)
		DeleteEvent(ctx context.Context, id int64) (int64, error)
		GetUserDayEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error)
		GetUserWeekEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error)
		GetUserMonthEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error)
	}

	StorageConnection interface {
		PingContext(context.Context) error
	}
)

type EventServiceServer struct {
	pb.UnimplementedEventServiceServer

	eventUseCase EventUseCase
	storageConn  StorageConnection
}

func NewEventServiceServer(eventUseCase EventUseCase, storageConn StorageConnection) *EventServiceServer {
	return &EventServiceServer{
		eventUseCase: eventUseCase,
		storageConn:  storageConn,
	}
}

func (es *EventServiceServer) GetEventByID(ctx context.Context, r *pb.GetEventByIDRequest) (*pb.GetEventByIDResponse, error) {
	e, err := es.eventUseCase.GetEventByID(ctx, r.Id)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "event not found")
	}
	if err != nil {
		return nil, err
	}

	return &pb.GetEventByIDResponse{Event: ToEvent(e)}, nil
}

func (es *EventServiceServer) CreateEvent(ctx context.Context, r *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	insertedID, err := es.eventUseCase.CreateEvent(ctx, FromEvent(r.Event))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, "date %s already busy", r.Event.StartDate.AsTime())
	}
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{InsertedID: insertedID}, nil
}

func (es *EventServiceServer) UpdateEvent(ctx context.Context, r *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	affected, err := es.eventUseCase.UpdateEvent(ctx, r.Id, FromEvent(r.Event))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, "date %s already busy", r.Event.StartDate.AsTime())
	}
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{Affected: affected}, nil
}

func (es *EventServiceServer) DeleteEvent(ctx context.Context, r *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	affected, err := es.eventUseCase.DeleteEvent(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteEventResponse{Affected: affected}, nil
}

func (es *EventServiceServer) GetUserDayEvents(ctx context.Context, r *pb.UserPeriodEventRequest) (*pb.EventListResponse, error) {
	events, err := es.eventUseCase.GetUserDayEvents(ctx, r.UserID, r.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.EventListResponse{Events: ToEventSlice(events)}, nil
}

func (es *EventServiceServer) GetUserWeekEvents(ctx context.Context, r *pb.UserPeriodEventRequest) (*pb.EventListResponse, error) {
	events, err := es.eventUseCase.GetUserWeekEvents(ctx, r.UserID, r.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.EventListResponse{Events: ToEventSlice(events)}, nil
}

func (es *EventServiceServer) GetUserMonthEvents(ctx context.Context, r *pb.UserPeriodEventRequest) (*pb.EventListResponse, error) {
	events, err := es.eventUseCase.GetUserMonthEvents(ctx, r.UserID, r.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.EventListResponse{Events: ToEventSlice(events)}, nil
}

func (es *EventServiceServer) Health(ctx context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := es.storageConn.PingContext(ctx); err != nil {
		return nil, ErrServiceUnavailable
	}

	return &pb.HealthResponse{Status: "alive"}, nil
}

func ToEvent(e model.Event) *pb.Event {
	return &pb.Event{
		Id:               e.ID,
		UserID:           e.UserID,
		Title:            e.Title,
		Description:      e.Description,
		StartDate:        timestamppb.New(e.StartDate),
		EndDate:          timestamppb.New(e.EndDate),
		NotificationDate: timestamppb.New(e.NotificationDate),
	}
}

func FromEvent(e *pb.Event) model.Event {
	return model.Event{
		ID:               e.Id,
		UserID:           e.UserID,
		Title:            e.Title,
		Description:      e.Description,
		StartDate:        e.StartDate.AsTime(),
		EndDate:          e.EndDate.AsTime(),
		NotificationDate: e.NotificationDate.AsTime(),
	}
}

func ToEventSlice(events []model.Event) []*pb.Event {
	pbEvents := make([]*pb.Event, 0, len(events))

	for _, e := range events {
		pbEvents = append(pbEvents, ToEvent(e))
	}

	return pbEvents
}
