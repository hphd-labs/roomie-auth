package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test Get log from request
func TestGetLog(t *testing.T) {
	t.Skip("This test is broken by powermux.",
		"Issue: https://github.com/andrewburian/powermux/issues/35")
	baseLog := logrus.NewEntry(logrus.StandardLogger()).WithField("test", "test")

	mid := &LoggerMiddleware{
		baseLog,
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Host = "test.com"
	req.RemoteAddr = "127.0.0.1:6969"

	mid.ServeHTTPMiddleware(nil, req, func(_ http.ResponseWriter, r *http.Request) {
		log := GetLog(r)
		if log == nil {
			t.Fatal("Nil entry returned")
		}
		if val := log.Data["test"]; val != "test" {
			t.Error("Missing base field", val)
		}
		if val := log.Data["path"]; val != req.URL.Path {
			t.Error("Mising path", val)
		}
		if val := log.Data["method"]; val != req.Method {
			t.Error("Missing method", val)
		}
		if val := log.Data["remote"]; val != req.RemoteAddr {
			t.Error("Missing remote", val)
		}

	})
}

// Test get a new log from a request without one
func TestGetLog2(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	log := GetLog(req)

	if log == nil {
		t.Fatal("Should have returned a log")
	}

	if len(log.Data) > 0 {
		t.Error("Fields set in what should be an empty entry")
	}
}
