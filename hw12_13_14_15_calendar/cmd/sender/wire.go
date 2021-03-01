// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/sender"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/factory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func setup(*config.Config) (*sender.Sender, func(), error) {
	panic(wire.Build(
		wire.Bind(new(sender.EventUseCase), new(*calendar.EventUseCase)),
		wire.Bind(new(sender.Queue), new(*rabbitmq.Rabbit)),
		sqlstorage.DatabaseProvider,
		factory.CreateEventRepository,
		calendar.NewEventUseCase,
		rabbitmq.NewRabbit,
		sender.NewSender,
	))
}
