package model

import (
	"testing"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestToEvent(t *testing.T) {
	se := storage.Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        time.Now(),
		EndDate:          time.Now(),
		NotificationDate: time.Now(),
	}

	expected := Event{
		ID:               int64(se.ID),
		Title:            se.Title,
		Description:      se.Description,
		UserID:           int64(se.UserID),
		StartDate:        se.StartDate,
		EndDate:          se.EndDate,
		NotificationDate: se.NotificationDate,
	}

	require.Equal(t, expected, ToEvent(se))
}

func TestFromEvent(t *testing.T) {
	e := Event{
		ID:               1,
		Title:            "title",
		Description:      "description",
		UserID:           1,
		StartDate:        time.Now(),
		EndDate:          time.Now(),
		NotificationDate: time.Now(),
	}

	expected := storage.Event{
		ID:               storage.EventID(e.ID),
		Title:            e.Title,
		Description:      e.Description,
		UserID:           storage.UserID(e.UserID),
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		NotificationDate: e.NotificationDate,
	}

	require.Equal(t, expected, FromEvent(e))
}

func TestToEventSlice(t *testing.T) {
	se := []storage.Event{
		{
			ID:               1,
			Title:            "title",
			Description:      "description",
			UserID:           1,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		},
		{
			ID:               2,
			Title:            "title2",
			Description:      "description2",
			UserID:           2,
			StartDate:        time.Now(),
			EndDate:          time.Now(),
			NotificationDate: time.Now(),
		},
	}

	expected := []Event{
		{
			ID:               int64(se[0].ID),
			Title:            se[0].Title,
			Description:      se[0].Description,
			UserID:           int64(se[0].UserID),
			StartDate:        se[0].StartDate,
			EndDate:          se[0].EndDate,
			NotificationDate: se[0].NotificationDate,
		},
		{
			ID:               int64(se[1].ID),
			Title:            se[1].Title,
			Description:      se[1].Description,
			UserID:           int64(se[1].UserID),
			StartDate:        se[1].StartDate,
			EndDate:          se[1].EndDate,
			NotificationDate: se[1].NotificationDate,
		},
	}

	require.Equal(t, expected, ToEventSlice(se))
}
