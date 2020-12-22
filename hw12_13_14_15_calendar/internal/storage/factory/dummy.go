package factory

import "context"

type dummyConnection struct{}

func (dc dummyConnection) PingContext(_ context.Context) error {
	return nil
}
