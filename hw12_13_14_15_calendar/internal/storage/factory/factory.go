package factory

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

var ErrUnexpectedStorage = errors.New("unexpected storage")

func CreateEventRepository(cfg *config.Config, db *sqlx.DB) (calendar.EventRepository, error) {
	switch cfg.StorageType {
	case config.InMemoryStorage:
		return memorystorage.NewEventStorage(), nil
	case config.SQLStorage:
		return sqlstorage.NewEventStorage(db), nil
	}

	return nil, ErrUnexpectedStorage
}
