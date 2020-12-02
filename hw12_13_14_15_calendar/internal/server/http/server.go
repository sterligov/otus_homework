package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	httpServer http.Server
}

func NewServer(cfg *config.Config, h http.Handler) (*Server, error) {
	rt, err := time.ParseDuration(cfg.HTTP.ReadTimeout)
	if err != nil {
		return nil, err
	}

	wt, err := time.ParseDuration(cfg.HTTP.ReadTimeout)
	if err != nil {
		return nil, err
	}

	ht, err := time.ParseDuration(cfg.HTTP.HandlerTimeout)
	if err != nil {
		return nil, err
	}

	server := &Server{
		httpServer: http.Server{
			Addr:         cfg.HTTP.Addr,
			ReadTimeout:  rt,
			WriteTimeout: wt,
			Handler:      http.TimeoutHandler(h, ht, "request timeout"),
		},
	}

	return server, nil
}

func (s *Server) Start() error {
	logger.Infof("Start server...")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Infof("Stop server...")

	return s.httpServer.Shutdown(ctx)
}
