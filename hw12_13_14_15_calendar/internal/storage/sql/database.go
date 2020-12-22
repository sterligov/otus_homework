package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
)

func NewDatabase(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, cfg.Database.Driver, cfg.Database.Addr)
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

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	db, err := NewDatabase(ctx, cfg)
	if err != nil {
		return nil, dbClose, err
	}

	dbClose = func() {
		if err := db.Close(); err != nil {
			logrus.Warnf("database close failed: %s", err)
		}
	}

	return db, dbClose, nil
}
