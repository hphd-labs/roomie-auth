package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"net/http"
)

type LoggerMiddleware struct {
	Base *logrus.Entry
}

type logMiddlewareCtxKey string

const (
	LOG_KEY = logMiddlewareCtxKey("log")
)

// Injects a log entry into the request for use by any handler
func (m *LoggerMiddleware) ServeHTTPMiddleware(w http.ResponseWriter, r *http.Request, n func(http.ResponseWriter, *http.Request)) {
	ctx := r.Context()
	e := m.Base.WithFields(map[string]interface{}{
		"path":   powermux.RequestPath(r),
		"method": r.Method,
		"remote": r.RemoteAddr,
		"host":   r.Host,
	})
	ctx = context.WithValue(ctx, LOG_KEY, e)
	e.Debug("Recieved request")
	n(w, r.WithContext(ctx))
}

// Retrieves a log entry from the context for use, or makes a new one from scratch
func GetLog(r *http.Request) *logrus.Entry {
	log, ok := r.Context().Value(LOG_KEY).(*logrus.Entry)
	if ok && log != nil {
		return log
	}

	return logrus.NewEntry(logrus.StandardLogger())
}
