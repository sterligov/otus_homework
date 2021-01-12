package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
)

const (
	dateLayout        = "2006-01-02 15:04:05"
	mysqlUniqueErrNum = 1062
)

type EventStorage struct {
	db *sqlx.DB
}

func NewEventStorage(db *sqlx.DB) *EventStorage {
	return &EventStorage{db: db}
}

func (es *EventStorage) GetEventByID(ctx context.Context, id storage.EventID) (storage.Event, error) {
	query := `
SELECT
	*
FROM 
    event
WHERE 
    id = ?`

	var event storage.Event

	row := es.db.QueryRowxContext(ctx, query, id)
	if err := row.StructScan(&event); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Event{}, storage.ErrNotFound
		}

		return storage.Event{}, err
	}

	return event, nil
}

func (es *EventStorage) CreateEvent(ctx context.Context, e storage.Event) (storage.EventID, error) {
	query := `
INSERT INTO event(
    title,
    description,
    user_id,
    start_date,
    end_date,
    notification_date
) VALUES (
    :title,
    :description,
    :user_id,
    :start_date,
    :end_date,
    :notification_date
)`

	res, err := es.db.NamedExecContext(ctx, query, &e)
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == mysqlUniqueErrNum {
			return 0, storage.ErrDateBusy
		}

		return 0, fmt.Errorf("create event failed: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id failed: %w", err)
	}

	return storage.EventID(lastID), nil
}

func (es *EventStorage) UpdateEvent(ctx context.Context, e storage.Event) (int64, error) {
	query := `
UPDATE
	event
SET
	title = :title,
	description = :description,
	user_id = :user_id,
	start_date = :start_date,
	end_date = :end_date,
	notification_date = :notification_date
WHERE
	id = :id`

	res, err := es.db.NamedExecContext(ctx, query, &e)
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok && me.Number == mysqlUniqueErrNum {
			return 0, storage.ErrDateBusy
		}

		return 0, fmt.Errorf("update event failed: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get affected rows failed: %w", err)
	}

	return affected, nil
}

func (es *EventStorage) DeleteEvent(ctx context.Context, id storage.EventID) (int64, error) {
	query := `DELETE FROM event WHERE id = ?`

	res, err := es.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("delete event failed: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get affected rows failed: %w", err)
	}

	return affected, nil
}

func (es *EventStorage) GetUserEventsByPeriod(
	ctx context.Context,
	uid storage.UserID,
	startDate, endDate time.Time,
) ([]storage.Event, error) {
	query := `
SELECT
	*
FROM
	event
WHERE
    user_id = ? AND start_date BETWEEN ? AND ?
ORDER BY
	start_date`

	sdate := startDate.Format(dateLayout)
	edate := endDate.Format(dateLayout)

	rows, err := es.db.QueryxContext(ctx, query, uid, sdate, edate)
	if err != nil {
		return nil, fmt.Errorf("fetching events failed: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Errorf("rows close failed: %s", err)
		}
	}()

	var (
		events []storage.Event
		event  storage.Event
	)

	for rows.Next() {
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("scan event failed: %w", err)
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}

	return events, nil
}

func (es *EventStorage) GetEventsByNotificationDatePeriod(
	ctx context.Context,
	startDate, endDate time.Time,
) ([]storage.Event, error) {
	query := `
SELECT
	*
FROM
	event
WHERE
    notification_date BETWEEN ? AND ?
ORDER BY
	notification_date`

	sdate := startDate.Format(dateLayout)
	edate := endDate.Format(dateLayout)

	rows, err := es.db.QueryxContext(ctx, query, sdate, edate)
	if err != nil {
		return nil, fmt.Errorf("fetching events failed: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Warnf("rows close failed: %s", err)
		}
	}()

	var (
		events []storage.Event
		event  storage.Event
	)

	for rows.Next() {
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("scan event failed: %w", err)
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}

	return events, nil
}

func (es *EventStorage) DeleteEventsBeforeDate(ctx context.Context, date time.Time) (int64, error) {
	query := `DELETE FROM event WHERE start_date <= ?`

	res, err := es.db.ExecContext(ctx, query, date)
	if err != nil {
		return 0, fmt.Errorf("delete events failed: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get affected rows failed: %w", err)
	}

	return affected, nil
}
