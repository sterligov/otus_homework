package sqlstorage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
)

func NewDatabase(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(cfg.Database.Driver, cfg.Database.Addr)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return db, nil
}

func DatabaseProvider(cfg *config.Config) (*sqlx.DB, func(), error) {
	dbClose := func() {}
	if cfg.StorageType != config.SQLStorage {
		return nil, dbClose, nil
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.MaxConnLifetime)

	dbClose = func() {
		if err := db.Close(); err != nil {
			logrus.Warnf("database close failed: %s", err)
		}
	}

	return db, dbClose, nil
}
