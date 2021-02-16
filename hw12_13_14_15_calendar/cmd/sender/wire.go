// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/sender"
)

func setup(*config.Config) (*sender.Sender, func(), error) {
	panic(wire.Build(
		wire.Bind(new(sender.Queue), new(*rabbitmq.Rabbit)),
		rabbitmq.NewRabbit,
		sender.NewSender,
	))
}
