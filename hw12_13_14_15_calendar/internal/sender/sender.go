package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/streadway/amqp"
)

type (
	Event struct {
		ID     int64
		UserID int64
		Title  string
		Date   time.Time
	}

	Queue interface {
		Consume(context.Context, rabbitmq.Handler) error
		Shutdown() error
	}

	EventUseCase interface {
		Notify(ctx context.Context, id int64) error
	}

	Sender struct {
		queue        Queue
		eventUseCase EventUseCase
	}
)

func NewSender(queue Queue, eventUseCase EventUseCase) *Sender {
	return &Sender{
		queue:        queue,
		eventUseCase: eventUseCase,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	logrus.Infof("Start sender...")

	return s.queue.Consume(ctx, s.Handle)
}

func (s *Sender) Shutdown() error {
	logrus.Infof("Stop sender...")

	return s.queue.Shutdown()
}

func (s *Sender) Handle(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			logrus.Infof("received message from queue")

			e := &Event{}

			if err := json.Unmarshal(msg.Body, e); err != nil {
				logrus.WithError(err).Errorf("unmarshal failed")
				return
			}

			if err := msg.Ack(false); err != nil {
				logrus.WithError(err).Error("ack failed")
			}

			if err := s.eventUseCase.Notify(ctx, e.ID); err != nil {
				logrus.
					WithField("event", e).
					WithError(err).
					Error("notify failed")
			}
		}
	}
}
