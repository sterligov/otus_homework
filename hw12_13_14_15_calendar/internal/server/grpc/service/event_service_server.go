//go:generate protoc -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis -I ./../../../../ --go_out ./../pb --go-grpc_out ./../pb ./../../../../api/event_service.proto
//go:generate protoc -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis -I ./../../../../ --grpc-gateway_out ./../pb --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true ./../../../../api/event_service.proto
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

func (es *EventServiceServer) GetEventByID(ctx context.Context, eid *pb.EventID) (*pb.Event, error) {
	e, err := es.eventUseCase.GetEventByID(ctx, eid.Id)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "event not found")
	}
	if err != nil {
		return nil, err
	}

	return ToEvent(e), nil
}

func (es *EventServiceServer) CreateEvent(ctx context.Context, pbEvent *pb.Event) (*pb.Inserted, error) {
	insertedID, err := es.eventUseCase.CreateEvent(ctx, FromEvent(pbEvent))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, "date %s already busy", pbEvent.StartDate.AsTime())
	}
	if err != nil {
		return nil, err
	}

	return &pb.Inserted{InsertedID: insertedID}, nil
}

func (es *EventServiceServer) UpdateEvent(ctx context.Context, e *pb.Event) (*pb.Affected, error) {
	affected, err := es.eventUseCase.UpdateEvent(ctx, e.Id, FromEvent(e))
	if errors.Is(err, storage.ErrDateBusy) {
		return nil, status.Errorf(codes.InvalidArgument, "date %s already busy", e.StartDate.AsTime())
	}
	if err != nil {
		return nil, err
	}

	return &pb.Affected{Affected: affected}, nil
}

func (es *EventServiceServer) DeleteEvent(ctx context.Context, eid *pb.EventID) (*pb.Affected, error) {
	affected, err := es.eventUseCase.DeleteEvent(ctx, eid.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Affected{Affected: affected}, nil
}

func (es *EventServiceServer) GetUserDayEvents(ctx context.Context, pr *pb.UserPeriodEventRequest) (*pb.Events, error) {
	events, err := es.eventUseCase.GetUserDayEvents(ctx, pr.UserID, pr.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return ToEvents(events), nil
}

func (es *EventServiceServer) GetUserWeekEvents(ctx context.Context, pr *pb.UserPeriodEventRequest) (*pb.Events, error) {
	events, err := es.eventUseCase.GetUserWeekEvents(ctx, pr.UserID, pr.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return ToEvents(events), nil
}

func (es *EventServiceServer) GetUserMonthEvents(ctx context.Context, pr *pb.UserPeriodEventRequest) (*pb.Events, error) {
	events, err := es.eventUseCase.GetUserMonthEvents(ctx, pr.UserID, pr.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return ToEvents(events), nil
}

func (es *EventServiceServer) Health(ctx context.Context, _ *pb.Empty) (*pb.HealthResponse, error) {
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

func ToEvents(events []model.Event) *pb.Events {
	pbEvents := &pb.Events{}

	for _, e := range events {
		pbEvents.Events = append(pbEvents.Events, ToEvent(e))
	}

	return pbEvents
}
