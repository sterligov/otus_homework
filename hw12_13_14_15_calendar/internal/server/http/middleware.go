package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseWriterDecorator struct {
	http.ResponseWriter

	status int
}

func newResponseWriterDecorator(rw http.ResponseWriter) *responseWriterDecorator {
	return &responseWriterDecorator{
		ResponseWriter: rw,
	}
}

func (rw *responseWriterDecorator) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wd, ok := w.(*responseWriterDecorator)
		if !ok {
			wd = newResponseWriterDecorator(w)
		}

		t := time.Now()
		next.ServeHTTP(wd, r)
		latency := fmt.Sprintf("%dms", time.Since(t).Milliseconds())

		logrus.Infof(
			"%s %s %s %s %d %s %s",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			r.Proto,
			wd.status,
			latency,
			r.UserAgent(),
		)
	})
}

func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
