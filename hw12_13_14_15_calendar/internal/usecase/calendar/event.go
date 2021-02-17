package calendar

import (
	"context"
	"time"

	"github.com/jinzhu/now"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/model"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type EventRepository interface {
	GetEventByID(ctx context.Context, id storage.EventID) (storage.Event, error)
	CreateEvent(ctx context.Context, event storage.Event) (storage.EventID, error)
	UpdateEvent(ctx context.Context, event storage.Event) (int64, error)
	DeleteEvent(ctx context.Context, id storage.EventID) (int64, error)
	GetEventsByNotificationDatePeriod(ctx context.Context, start, end time.Time) ([]storage.Event, error)
	DeleteNotifiedEventsBeforeDate(ctx context.Context, date time.Time) (int64, error)
	GetUserEventsByPeriod(ctx context.Context, uid storage.UserID, start, end time.Time) ([]storage.Event, error)
	UpdateIsNotified(ctx context.Context, id storage.EventID, isNotified byte) error
}

type EventUseCase struct {
	eventRepository EventRepository
}

func NewEventUseCase(eventRepository EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepository: eventRepository,
	}
}

func (eu *EventUseCase) GetEventByID(ctx context.Context, id int64) (model.Event, error) {
	e, err := eu.eventRepository.GetEventByID(ctx, storage.EventID(id))
	if err != nil {
		return model.Event{}, err
	}

	return model.ToEvent(e), nil
}

func (eu *EventUseCase) CreateEvent(ctx context.Context, e model.Event) (int64, error) {
	insertedID, err := eu.eventRepository.CreateEvent(ctx, model.FromEvent(e))
	if err != nil {
		return 0, err
	}

	return int64(insertedID), nil
}

func (eu *EventUseCase) UpdateEvent(ctx context.Context, id int64, e model.Event) (int64, error) {
	e.ID = id

	affected, err := eu.eventRepository.UpdateEvent(ctx, model.FromEvent(e))
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (eu *EventUseCase) DeleteEvent(ctx context.Context, id int64) (int64, error) {
	return eu.eventRepository.DeleteEvent(ctx, storage.EventID(id))
}

func (eu *EventUseCase) GetUserDayEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error) {
	start := now.With(date).BeginningOfDay()
	end := now.With(date).EndOfDay()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return model.ToEventSlice(events), nil
}

func (eu *EventUseCase) GetUserWeekEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error) {
	start := now.With(date).BeginningOfWeek()
	end := now.With(date).EndOfWeek()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return model.ToEventSlice(events), nil
}

func (eu *EventUseCase) GetUserMonthEvents(ctx context.Context, uid int64, date time.Time) ([]model.Event, error) {
	start := now.With(date).BeginningOfMonth()
	end := now.With(date).EndOfMonth()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return model.ToEventSlice(events), nil
}

func (eu *EventUseCase) DeleteNotifiedEventsBeforeDate(ctx context.Context, date time.Time) (int64, error) {
	affected, err := eu.eventRepository.DeleteNotifiedEventsBeforeDate(ctx, date)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (eu *EventUseCase) GetEventsByNotificationDatePeriod(ctx context.Context, start, end time.Time) ([]model.Event, error) {
	events, err := eu.eventRepository.GetEventsByNotificationDatePeriod(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return model.ToEventSlice(events), nil
}

func (eu *EventUseCase) Notify(ctx context.Context, id int64) error {
	return eu.eventRepository.UpdateIsNotified(ctx, storage.EventID(id), 1)
}
