// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/service"
	internalhttp "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/factory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func setup(*config.Config) (*server.Server, func(), error) {
	panic(wire.Build(
		wire.Bind(new(service.EventUseCase), new(*calendar.EventUseCase)),
		wire.Bind(new(pb.EventServiceServer), new(*service.EventServiceServer)),
		factory.GetStorageConnection,
		sqlstorage.DatabaseProvider,
		factory.CreateEventRepository,
		calendar.NewEventUseCase,
		service.NewEventServiceServer,
		internalhttp.NewHandler,
		internalhttp.NewServer,
		internalgrpc.NewServer,
		server.NewServer,
	))
}
