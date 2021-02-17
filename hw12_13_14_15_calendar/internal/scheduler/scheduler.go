package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/model"
)

type (
	Event struct {
		ID     int64
		UserID int64
		Title  string
		Date   time.Time
	}

	EventAlias Event

	EventUseCase interface {
		UpdateEvent(ctx context.Context, id int64, e model.Event) (int64, error)
		DeleteNotifiedEventsBeforeDate(ctx context.Context, date time.Time) (int64, error)
		GetEventsByNotificationDatePeriod(ctx context.Context, start, end time.Time) ([]model.Event, error)
	}

	Queue interface {
		Publish(context.Context, json.Marshaler) error
		Shutdown() error
	}

	Scheduler struct {
		queue        Queue
		frequency    time.Duration
		eventUseCase EventUseCase
	}
)

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(EventAlias(e))
}

func NewScheduler(
	cfg *config.Config,
	queue Queue,
	eventUseCase EventUseCase,
) *Scheduler {
	return &Scheduler{
		queue:        queue,
		frequency:    cfg.EventScanFreq,
		eventUseCase: eventUseCase,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	logrus.Infof("Start scheduler...")

	ticker := time.NewTicker(s.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Infof("Stop scheduler...")

			return ctx.Err()
		case <-ticker.C:
			go func() {
				s.sendNotifications(ctx)
				s.deleteOldNotifiedEvents(ctx)
			}()
		}
	}
}

func (s *Scheduler) Shutdown() error {
	logrus.Info("Stop scheduler...")

	return s.queue.Shutdown()
}

func (s *Scheduler) sendNotifications(ctx context.Context) {
	edate := time.Now()
	sdate := time.Now().Add(-s.frequency + time.Second)

	events, err := s.eventUseCase.GetEventsByNotificationDatePeriod(ctx, sdate, edate)
	if err != nil {
		logrus.WithError(err).Error("get events by period failed")
		return
	}

	var published int

	for _, e := range events {
		if err := s.queue.Publish(ctx, ToEvent(e)); err != nil {
			logrus.
				WithError(err).
				WithField("event", e).
				Error("publish as json failed")
			continue
		}

		published++
	}

	logrus.Infof("%d events were published successfully, %d errors", published, len(events)-published)
}

func (s *Scheduler) deleteOldNotifiedEvents(ctx context.Context) {
	t := time.Now().AddDate(-1, 0, 0)
	log := logrus.WithField("date", t.String())

	affected, err := s.eventUseCase.DeleteNotifiedEventsBeforeDate(ctx, t)
	if err != nil {
		logrus.WithError(err).Error("delete old events failed")
		return
	}

	log.Infof("delete old events: affected rows %d", affected)
}

func ToEvent(e model.Event) Event {
	return Event{
		ID:     e.ID,
		Title:  e.Title,
		UserID: e.UserID,
		Date:   e.StartDate,
	}
}
