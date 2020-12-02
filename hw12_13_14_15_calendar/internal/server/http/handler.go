package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/usecase/calendar"
)

type EventUseCase interface {
	GetEventByID(ctx context.Context, id int64) (calendar.Event, error)
	CreateEvent(ctx context.Context, e calendar.Event) (calendar.Event, error)
	UpdateEvent(ctx context.Context, id int64, e calendar.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetUserDayEvents(ctx context.Context, uid int64, date time.Time) ([]calendar.Event, error)
	GetUserWeekEvents(ctx context.Context, uid int64, date time.Time) ([]calendar.Event, error)
	GetUserMonthEvents(ctx context.Context, uid int64, date time.Time) ([]calendar.Event, error)
}

type handler struct {
	eventUseCase EventUseCase
}

func NewHandler(eu EventUseCase) http.Handler {
	h := handler{
		eventUseCase: eu,
	}

	healthHandler := HeadersMiddleware(http.HandlerFunc(h.Health))
	healthHandler = LoggingMiddleware(healthHandler)
	healthHandler = RecoverMiddleware(healthHandler)

	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler)

	return mux
}

func (h *handler) Health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte(`{"status":"alive"}`)); err != nil {
		logger.Named("health check").Errorf("write answer error: %s", err.Error())
	}
}
