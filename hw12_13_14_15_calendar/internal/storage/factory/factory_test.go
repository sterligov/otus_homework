package factory

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

func TestCreateEventRepository(t *testing.T) {
	tests := []struct {
		config  config.Config
		repType calendar.EventRepository
		err     error
	}{
		{
			config:  config.Config{StorageType: config.SQLStorage},
			repType: &sqlstorage.EventStorage{},
			err:     nil,
		},
		{
			config:  config.Config{StorageType: config.InMemoryStorage},
			repType: &memorystorage.EventStorage{},
			err:     nil,
		},
		{
			config:  config.Config{StorageType: "unexpected storage"},
			repType: nil,
			err:     ErrUnexpectedStorage,
		},
	}

	for _, tst := range tests {
		t.Run(tst.config.StorageType, func(t *testing.T) {
			rep, err := CreateEventRepository(&tst.config, nil)
			require.Equal(t, tst.err, err)
			require.IsType(t, tst.repType, rep)
		})
	}
}
