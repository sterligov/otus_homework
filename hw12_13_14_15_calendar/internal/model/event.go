package model

import (
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type Event struct {
	ID               int64
	Title            string
	Description      string
	UserID           int64
	StartDate        time.Time
	EndDate          time.Time
	NotificationDate time.Time
	IsNotified       byte
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
		IsNotified:       e.IsNotified,
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
		IsNotified:       e.IsNotified,
	}
}

func ToEventSlice(storEvents []storage.Event) []Event {
	events := make([]Event, 0, len(storEvents))

	for _, e := range storEvents {
		events = append(events, ToEvent(e))
	}

	return events
}
