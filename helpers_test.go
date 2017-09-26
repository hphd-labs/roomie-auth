package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAcceptsJSON(t *testing.T) {
	// test accept case
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")

	if !AcceptsJSON(req) {
		t.Error("Should report JSON accepted")
	}
}

func TestAcceptsJSON2(t *testing.T) {
	// test accepts any application/*
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/*")

	if !AcceptsJSON(req) {
		t.Error("Should report JSON accepted with application/*")
	}
}

func TestAcceptsJSON3(t *testing.T) {
	// test body is json
	body := bytes.NewBufferString("{}")
	req := httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if !AcceptsJSON(req) {
		t.Error("Should report JSON accepted with body")
	}
}

func TestAcceptsJSON4(t *testing.T) {
	// test reject unknown body type
	body := bytes.NewBufferString("hi")
	req := httptest.NewRequest(http.MethodGet, "/", body)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Accept", "application/json")

	if AcceptsJSON(req) {
		t.Error("Should report body not acceptable")
	}
}

func TestRenderJSON(t *testing.T) {
	res := httptest.NewRecorder()

	data := &struct {
		Message string
	}{
		Message: "hi",
	}

	err := RenderJSON(res, data)
	if err != nil {
		t.Fatal(err)
	}

	if val := res.Header().Get("Content-Type"); val != "application/json" {
		t.Error("Invalid content type", val)
	}
}
