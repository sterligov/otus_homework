package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type EventStorage struct {
	mu     sync.RWMutex
	bucket map[storage.EventID]storage.Event
	lastID storage.EventID
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		bucket: make(map[storage.EventID]storage.Event),
	}
}

func (es *EventStorage) GetEventByID(_ context.Context, id storage.EventID) (storage.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	if e, ok := es.bucket[id]; ok {
		return e, nil
	}

	return storage.Event{}, storage.ErrNotFound
}

func (es *EventStorage) CreateEvent(_ context.Context, event storage.Event) (storage.EventID, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	for _, e := range es.bucket {
		if e.StartDate.Equal(event.StartDate) && e.UserID == event.UserID {
			return 0, storage.ErrDateBusy
		}
	}

	es.lastID++
	event.ID = es.lastID
	es.bucket[es.lastID] = event

	return es.lastID, nil
}

func (es *EventStorage) UpdateEvent(_ context.Context, event storage.Event) (int64, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, ok := es.bucket[event.ID]; !ok {
		return 0, nil
	}

	for _, e := range es.bucket {
		if e.ID != event.ID && e.StartDate.Equal(event.StartDate) && e.UserID == event.UserID {
			return 0, storage.ErrDateBusy
		}
	}

	es.bucket[event.ID] = event

	return 1, nil
}

func (es *EventStorage) UpdateIsNotified(_ context.Context, id storage.EventID, isNotified byte) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	for k, e := range es.bucket {
		if e.ID == id {
			e.IsNotified = isNotified
			es.bucket[k] = e
			return nil
		}
	}

	return nil
}

func (es *EventStorage) DeleteEvent(_ context.Context, id storage.EventID) (int64, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	if _, ok := es.bucket[id]; !ok {
		return 0, nil
	}

	delete(es.bucket, id)

	return 1, nil
}

func (es *EventStorage) GetUserEventsByPeriod(
	_ context.Context,
	uid storage.UserID,
	startDate, endDate time.Time,
) ([]storage.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var events []storage.Event

	for _, e := range es.bucket {
		if e.UserID == uid && startDate.Before(e.StartDate) && endDate.After(e.StartDate) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (es *EventStorage) GetEventsByNotificationDatePeriod(
	_ context.Context,
	startDate, endDate time.Time,
) ([]storage.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var events []storage.Event

	for _, e := range es.bucket {
		if startDate.Before(e.StartDate) && endDate.After(e.StartDate) {
			events = append(events, e)
		}
	}

	return events, nil
}

func (es *EventStorage) DeleteNotifiedEventsBeforeDate(_ context.Context, date time.Time) (int64, error) {
	es.mu.Lock()
	defer es.mu.Unlock()

	var deleted int64

	for k, e := range es.bucket {
		if date.After(e.StartDate) && e.IsNotified == 1 {
			delete(es.bucket, k)
			deleted++
		}
	}

	return deleted, nil
}
