package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	ErrPeerFromContext = status.Error(codes.Internal, "get peer from context failed")
	ErrInternalError   = status.Error(codes.Internal, "internal server error")
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		logrus.Error(ErrPeerFromContext)

		return resp, ErrPeerFromContext
	}

	t := time.Now()
	resp, err = handler(ctx, req)
	latency := fmt.Sprintf("%dms", time.Since(t).Milliseconds())

	logrus.WithError(err).Infof(
		"%s %s %s %s",
		status.Code(err),
		info.FullMethod,
		p.Addr,
		latency,
	)

	return resp, err
}

func ErrorInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = ErrInternalError
		}
	}()

	resp, err = handler(ctx, req)

	code := status.Code(err)
	if code == codes.Unknown || code == codes.Internal {
		return resp, ErrInternalError
	}

	return resp, err
}
