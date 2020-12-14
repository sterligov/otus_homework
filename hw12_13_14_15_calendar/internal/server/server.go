package server

import (
	internalgrpc "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/http"
)

type Server struct {
	GRPC *internalgrpc.Server
	HTTP *internalhttp.Server
}

func NewServer(grpcServer *internalgrpc.Server, httpServer *internalhttp.Server) *Server {
	return &Server{
		GRPC: grpcServer,
		HTTP: httpServer,
	}
}
