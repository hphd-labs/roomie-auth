package main

import (
	"context"
	"github.com/andrewburian/powermux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// PasswordAuthHandler deals with password based authentication
// attempts
type PasswordAuthHandler struct {
	DB           *AuthDatabase
	PasswordCost int
}

func (h *PasswordAuthHandler) Setup(r *powermux.Route) {
	// Register routes
	r.PostFunc(h.Login)

	// Set password hash cost
	if h.PasswordCost == 0 {
		h.PasswordCost = bcrypt.DefaultCost
	}
}

func (h *PasswordAuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// context from request
	ctx := r.Context()

	// JSON check
	if !AcceptsJSON(r) {
		http.Error(w, "JSON must be provided and received", http.StatusNotAcceptable)
		return
	}

	// Decode request body
	var authAttempt PasswordLogin

	if DecodeJSON(r.Body, &authAttempt) != nil {
		http.Error(w, "Failed to decode auth attempt", http.StatusBadRequest)
		return
	}

	// Create the user
	user := &User{
		Username: authAttempt.Username,
	}

	// Get the user's credentials
	err := h.DB.GetUserByName(ctx, user)
	if err != nil {
		switch err {
		case errUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// check user credentials
	err = bcrypt.CompareHashAndPassword(user.Credentials, []byte(authAttempt.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusForbidden)
		return
	}

	// check if the password needs upgraded
	// this can be done in the background as it doesn't affect this endpoint's return
	go h.UpgradeAuth(context.Background(), user, []byte(authAttempt.Password))

	// Cut user an auth token
	//TODO generate token
	token := user.ID

	// Render response
	resp := &PasswordAuthReply{
		Token: token,
	}

	err = RenderJSON(w, resp)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

// UpgradeAuth checks if the password's bcrypt cost is up to current standards, and rehashes it if it's not,
// overwriting the existing creds with a new higher cost version.
func (h *PasswordAuthHandler) UpgradeAuth(ctx context.Context, user *User, password []byte) {
	cost, err := bcrypt.Cost(user.Credentials)
	if err != nil {
		//TODO log this very loudly
		return
	}

	if cost >= h.PasswordCost {
		// Nothing to be done
		return
	}

	newCreds, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		//TODO log
		return
	}

	// new user object to hold the new creds
	updatedUser := &User{
		ID:          user.ID,
		Credentials: newCreds,
	}

	// Store new creds
	err = h.DB.SetUserCredentials(ctx, updatedUser)
	if err != nil {
		//TODO log
	}
}
