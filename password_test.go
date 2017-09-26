package main

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test successful login
func TestPasswordAuthHandler_Login(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB:           newMockDatabase(),
		PasswordCost: bcrypt.MinCost,
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

// Test Incorrect password
func TestPasswordAuthHandler_Login2(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB: newMockDatabase(),
	}

	body := bytes.NewBufferString(`{"username":"test","password":"notright"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/password", body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler.Login(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatal("Login should have failed", res.Body.String())
	}
}

// Test Nonexistant user
func TestPasswordAuthHandler_Login3(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB: newMockDatabase(),
	}

	body := bytes.NewBufferString(`{"username":"otherguy","password":"password"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/password", body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler.Login(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatal("Login should have failed", res.Body.String())
	}
}

// Test password needs upgraded
func TestPasswordAuthHandler_UpgradeAuth(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB:           newMockDatabase(),
		PasswordCost: bcrypt.DefaultCost,
	}

	mock := handler.DB.(*MockDatabase)
	pass := mock.passwords[0]

	handler.UpgradeAuth(context.Background(), pass, []byte("password"))

	if mock.calls["SetPassword"] != 1 {
		t.Fatal("Should have set password in the database")
	}
}

// Test password doesn't need upgraded
func TestPasswordAuthHandler_UpgradeAuth1(t *testing.T) {
	handler := &PasswordAuthHandler{
		DB:           newMockDatabase(),
		PasswordCost: bcrypt.MinCost,
	}

	mock := handler.DB.(*MockDatabase)
	pass := mock.passwords[0]

	handler.UpgradeAuth(context.Background(), pass, []byte("password"))

	if mock.calls["SetPassword"] != 0 {
		t.Fatal("Should not have set password in the database")
	}
}
