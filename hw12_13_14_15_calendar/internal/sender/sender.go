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

	Sender struct {
		queue Queue
	}
)

func NewSender(queue Queue) *Sender {
	return &Sender{
		queue: queue,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	logrus.Infof("Start sender...")

	return s.queue.Consume(ctx, Handle)
}

func (s *Sender) Shutdown() error {
	logrus.Infof("Stop sender...")

	return s.queue.Shutdown()
}

func Handle(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}

			e := &Event{}

			if err := json.Unmarshal(msg.Body, e); err != nil {
				logrus.Errorf("unmarshal failed: %v", err)
			}

			logrus.Infof("received message from queue: %v", e)

			if err := msg.Ack(false); err != nil {
				logrus.Errorf("%s: ack failed: %v", msg.Body, err)
			}
		}
	}
}
