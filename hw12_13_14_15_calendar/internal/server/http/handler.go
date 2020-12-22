package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc"
)

func NewHandler(cfg *config.Config) (http.Handler, error) {
	gw := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEventServiceHandlerFromEndpoint(context.Background(), gw, cfg.GRPC.Addr, opts)
	if err != nil {
		return nil, fmt.Errorf("register event service handler endpoint failed: %w", err)
	}

	mux := http.NewServeMux()
	handler := HeadersMiddleware(gw)
	handler = LoggingMiddleware(handler)
	mux.Handle("/", handler)

	return mux, nil
}
