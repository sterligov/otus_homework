// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/factory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func setup(*config.Config) (*scheduler.Scheduler, func(), error) {
	panic(wire.Build(
		wire.Bind(new(scheduler.Queue), new(*rabbitmq.Rabbit)),
		wire.Bind(new(scheduler.EventUseCase), new(*calendar.EventUseCase)),
		sqlstorage.DatabaseProvider,
		rabbitmq.NewRabbitConnection,
		factory.CreateEventRepository,
		calendar.NewEventUseCase,
		scheduler.NewScheduler,
	))
}
