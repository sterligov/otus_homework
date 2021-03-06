// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/service"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/factory"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

// Injectors from wire.go:

func setup(cfg *config.Config) (*server.Server, func(), error) {
	db, cleanup, err := sqlstorage.DatabaseProvider(cfg)
	if err != nil {
		return nil, nil, err
	}
	eventRepository, err := factory.CreateEventRepository(cfg, db)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	eventUseCase := calendar.NewEventUseCase(eventRepository)
	storageConnection := factory.GetStorageConnection(db)
	eventServiceServer := service.NewEventServiceServer(eventUseCase, storageConnection)
	grpcServer := grpc.NewServer(cfg, eventServiceServer)
	handler, err := internalhttp.NewHandler(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	internalhttpServer, err := internalhttp.NewServer(cfg, handler)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	serverServer := server.NewServer(grpcServer, internalhttpServer)
	return serverServer, func() {
		cleanup()
	}, nil
}
