package calendar

import (
	"context"
	"time"

	"github.com/jinzhu/now"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	ID               int64     `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	UserID           int64     `json:"user_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	NotificationDate time.Time `json:"notification_date"`
}

type EventRepository interface {
	GetEventByID(ctx context.Context, id storage.EventID) (storage.Event, error)
	CreateEvent(ctx context.Context, event storage.Event) (storage.EventID, error)
	UpdateEvent(ctx context.Context, event storage.Event) (int64, error)
	DeleteEvent(ctx context.Context, id storage.EventID) (int64, error)
	GetUserEventsByPeriod(ctx context.Context, uid storage.UserID, start, end time.Time) ([]storage.Event, error)
}

type EventUseCase struct {
	eventRepository EventRepository
}

func NewEventUseCase(eventRepository EventRepository) *EventUseCase {
	return &EventUseCase{
		eventRepository: eventRepository,
	}
}

func (eu *EventUseCase) GetEventByID(ctx context.Context, id int64) (Event, error) {
	e, err := eu.eventRepository.GetEventByID(ctx, storage.EventID(id))
	if err != nil {
		return Event{}, err
	}

	return ToEvent(e), nil
}

func (eu *EventUseCase) CreateEvent(ctx context.Context, e Event) (int64, error) {
	if e.NotificationDate.Equal(time.Time{}) {
		e.NotificationDate = e.StartDate
	}

	insertedID, err := eu.eventRepository.CreateEvent(ctx, FromEvent(e))
	if err != nil {
		return 0, err
	}

	return int64(insertedID), nil
}

func (eu *EventUseCase) UpdateEvent(ctx context.Context, id int64, e Event) (int64, error) {
	e.ID = id

	affected, err := eu.eventRepository.UpdateEvent(ctx, FromEvent(e))
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (eu *EventUseCase) DeleteEvent(ctx context.Context, id int64) (int64, error) {
	return eu.eventRepository.DeleteEvent(ctx, storage.EventID(id))
}

func (eu *EventUseCase) GetUserDayEvents(ctx context.Context, uid int64, date time.Time) ([]Event, error) {
	start := now.With(date).BeginningOfDay()
	end := now.With(date).EndOfDay()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return ToEventSlice(events), nil
}

func (eu *EventUseCase) GetUserWeekEvents(ctx context.Context, uid int64, date time.Time) ([]Event, error) {
	start := now.With(date).BeginningOfWeek()
	end := now.With(date).EndOfWeek()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return ToEventSlice(events), nil
}

func (eu *EventUseCase) GetUserMonthEvents(ctx context.Context, uid int64, date time.Time) ([]Event, error) {
	start := now.With(date).BeginningOfMonth()
	end := now.With(date).EndOfMonth()

	events, err := eu.eventRepository.GetUserEventsByPeriod(ctx, storage.UserID(uid), start, end)
	if err != nil {
		return nil, err
	}

	return ToEventSlice(events), nil
}

func ToEvent(e storage.Event) Event {
	return Event{
		ID:               int64(e.ID),
		UserID:           int64(e.UserID),
		Title:            e.Title,
		Description:      e.Description,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		NotificationDate: e.NotificationDate,
	}
}

func FromEvent(e Event) storage.Event {
	return storage.Event{
		ID:               storage.EventID(e.ID),
		UserID:           storage.UserID(e.UserID),
		Title:            e.Title,
		Description:      e.Description,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		NotificationDate: e.NotificationDate,
	}
}

func ToEventSlice(storEvents []storage.Event) []Event {
	events := make([]Event, 0, len(storEvents))

	for _, e := range storEvents {
		events = append(events, ToEvent(e))
	}

	return events
}
