package scheduler

import (
	"context"
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

	EventUseCase interface {
		DeleteEventsBeforeDate(ctx context.Context, date time.Time) (int64, error)
		GetEventsByNotificationDatePeriod(ctx context.Context, start, end time.Time) ([]model.Event, error)
	}

	Queue interface {
		PublishAsJSON(context.Context, interface{}) error
		Shutdown() error
	}

	Scheduler struct {
		queue        Queue
		frequency    time.Duration
		eventUseCase EventUseCase
	}
)

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

func (s *Scheduler) Run(ctx context.Context) {
	logrus.Infof("Start scheduler...")

	ticker := time.NewTicker(s.frequency)

	for {
		select {
		case <-ctx.Done():
			logrus.Infof("Stop scheduler...")

			return
		case <-ticker.C:
			go func() {
				s.sendNotifications(ctx)
				s.deleteOldEvents(ctx)
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
		logrus.WithError(err).Errorf("get events by period failed")
		return
	}

	var published int

	for _, e := range events {
		if err := s.queue.PublishAsJSON(ctx, ToEvent(e)); err != nil {
			logrus.
				WithError(err).
				WithField("event", e).
				Errorf("publish as json failed")
			continue
		}
		published++
	}

	logrus.Infof("%d events were published successfully, %d errors", published, len(events)-published)
}

func (s *Scheduler) deleteOldEvents(ctx context.Context) {
	t := time.Now().AddDate(-1, 0, 0)
	log := logrus.WithField("date", t.String())

	affected, err := s.eventUseCase.DeleteEventsBeforeDate(ctx, t)
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
