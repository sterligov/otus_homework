package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"

	"github.com/jmoiron/sqlx"
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
	if cfg.StorageType != config.SQLStorage {
		return nil, func() {}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := NewDatabase(ctx, cfg)
	if err != nil {
		return nil, func() {}, err
	}

	dbClose := func() {
		if err := db.Close(); err != nil {
			logger.Warnf("database close failed: %s", err)
		}
	}

	return db, dbClose, nil
}
