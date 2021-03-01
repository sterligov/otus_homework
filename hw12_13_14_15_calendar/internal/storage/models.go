package storage

import "time"

type (
	EventID int64
	UserID  int64
)

type Event struct {
	ID               EventID   `db:"id"`
	Title            string    `db:"title"`
	Description      string    `db:"description"`
	UserID           UserID    `db:"user_id"`
	StartDate        time.Time `db:"start_date"`
	EndDate          time.Time `db:"end_date"`
	NotificationDate time.Time `db:"notification_date"`
	IsNotified       byte      `db:"is_notified"`
}
