package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test successful login
func TestPasswordAuthHandler_Login(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB: newMockDatabase(),
	}

	body := bytes.NewBufferString(`{"username":"test","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/password", body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler.Login(res, req)

	if res.Code != http.StatusOK {
		t.Fatal("Login should have succeeded", res.Body.String())
	}

	if ct := res.Header().Get("Content-Type"); ct != "application/json" {
		t.Error("Non-JSON repsonse type", ct)
	}

	if handler.DB.(*MockDatabase).calls["GetUserByName"] != 1 {
		t.Error("Should have called the database once")
	}

	var reply PasswordAuthReply
	if err := json.NewDecoder(res.Body).Decode(&reply); err != nil {
		t.Fatal("Failed to decode response", err.Error())
	}

	if reply.Token == "" {
		t.Error("No token set")
	}
}
