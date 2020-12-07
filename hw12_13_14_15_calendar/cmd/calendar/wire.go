// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/factory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func setup(*config.Config) (*internalhttp.Server, func(), error) {
	panic(wire.Build(
		wire.Bind(new(internalhttp.EventUseCase), new(*calendar.EventUseCase)),
		sqlstorage.DatabaseProvider,
		factory.CreateEventRepository,
		calendar.NewEventUseCase,
		internalhttp.NewHandler,
		internalhttp.NewServer,
	))
}
