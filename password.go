package main

import (
	"github.com/andrewburian/powermux"
	"net/http"
)

// PasswordAuthHandler deals with password based authentication
// attempts
type PasswordAuthHandler struct {
}

func (h *PasswordAuthHandler) Setup(r *powermux.Route) {
	r.PostFunc(h.Login)
}

func (h *PasswordAuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Decode request body
	var authAttempt PasswordLogin

	if DecodeJSON(r.Body, &authAttempt) != nil {
		http.Error(w, "Failed to decode auth attempt", http.StatusBadRequest)
		return
	}

}
