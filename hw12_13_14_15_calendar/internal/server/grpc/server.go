package grpc

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	addr       string
}

func NewServer(cfg *config.Config, eventServer pb.EventServiceServer) *Server {
	chainInterceptor := grpc.ChainUnaryInterceptor(
		LoggingInterceptor,
		ErrorInterceptor,
	)
	grpcServer := grpc.NewServer(chainInterceptor)
	pb.RegisterEventServiceServer(grpcServer, eventServer)

	return &Server{
		grpcServer: grpcServer,
		addr:       cfg.GRPC.Addr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("start grpc server failed: %w", err)
	}

	logrus.Infof("Start grpc server...")

	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop() {
	logrus.Infof("Stop grpc server...")

	s.grpcServer.GracefulStop()
}
