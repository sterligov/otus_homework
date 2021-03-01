package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

var ErrMaxReconnectRetries = errors.New("exceeded number of reconnect retries")

type (
	Rabbit struct {
		addr                string
		queueName           string
		nHandlers           int
		maxReconnectRetries int
		closed              int32
		reconnectInterval   time.Duration
		conn                *amqp.Connection
		channel             *amqp.Channel
	}

	Handler func(context.Context, <-chan amqp.Delivery)
)

func NewRabbit(cfg *config.Config) *Rabbit {
	return &Rabbit{
		addr:                cfg.AMQP.ConnectionAddr,
		queueName:           cfg.AMQP.QueueName,
		nHandlers:           cfg.AMQP.HandlersNumber,
		maxReconnectRetries: cfg.AMQP.MaxReconnectRetries,
		reconnectInterval:   cfg.AMQP.ReconnectInterval,
	}
}

func NewRabbitConnection(cfg *config.Config) (*Rabbit, error) {
	r := NewRabbit(cfg)

	if err := r.reConnect(context.Background()); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Rabbit) Consume(ctx context.Context, handler Handler) error {
	var err error

	if err = r.reConnect(ctx); err != nil {
		return err
	}

	for {
		msgs, err := r.announceQueue()
		if err != nil {
			return err
		}

		for i := 0; i < r.nHandlers; i++ {
			go handler(ctx, msgs)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case amqpErr, ok := <-r.conn.NotifyClose(make(chan *amqp.Error)):
			if !ok || r.isClosed() {
				return nil
			}

			logrus.WithError(amqpErr).Warn("close connection")

			if err := r.reConnect(ctx); err != nil {
				return fmt.Errorf("reconnecting failed: %w", err)
			}
		}
	}
}

func (r *Rabbit) Publish(ctx context.Context, marshaler json.Marshaler) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case amqpErr, ok := <-r.conn.NotifyClose(make(chan *amqp.Error)):
		if !ok || r.isClosed() {
			return nil
		}

		logrus.WithError(amqpErr).Warn("close connection")

		if err := r.reConnect(ctx); err != nil {
			return fmt.Errorf("reconnecting failed: %w", err)
		}
	default:
	}

	data, err := marshaler.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal json failed: %w", err)
	}

	return r.channel.Publish(
		"",
		r.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
}

func (r *Rabbit) Shutdown() error {
	r.close()

	if err := r.channel.Cancel("", true); err != nil {
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	logrus.Infof("AMQP shutdown")

	return nil
}

func (r *Rabbit) connect() error {
	var err error

	r.conn, err = amqp.Dial(r.addr)
	if err != nil {
		return fmt.Errorf("amqp dial failed: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel failed: %w", err)
	}

	logrus.Info("successfully connect")

	return nil
}

func (r *Rabbit) announceQueue() (<-chan amqp.Delivery, error) {
	_, err := r.channel.QueueDeclare(r.queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("queue declare failed: %w", err)
	}

	msgs, err := r.channel.Consume(
		r.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume failed: %w", err)
	}

	return msgs, nil
}

func (r *Rabbit) reConnect(ctx context.Context) error {
	var retries int
	for {
		retries++

		logrus.Infof("reconnect attempt %d", retries)

		if retries > r.maxReconnectRetries {
			return ErrMaxReconnectRetries
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.reconnectInterval):
			if err := r.connect(); err != nil {
				logrus.WithError(err).Warn("reconnect failed")
				break
			}

			return nil
		}
	}
}

func (r *Rabbit) isClosed() bool {
	return atomic.LoadInt32(&r.closed) == 1
}

func (r *Rabbit) close() {
	atomic.StoreInt32(&r.closed, 1)
}
