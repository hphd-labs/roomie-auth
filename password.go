package main

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/andrewburian/powermux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// PasswordAuthHandler deals with password based authentication
// attempts
type PasswordAuthHandler struct {
	DB           AuthDatabase
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
		case ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, "Database error", http.StatusInternalServerError)
			GetLog(r).WithField("component", "database").Error(err)
		}
		return
	}

	// check the user has password auth enabled
	if user.Password == nil {
		http.Error(w, "Password authorization not allowed for this user", http.StatusForbidden)
		GetLog(r).WithField("user", user.ID).Debug("User attempted password auth when no password was enabled")
		return
	}

	// check user credentials
	err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(authAttempt.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// check if the password needs upgraded
	// this can be done in the background as it doesn't affect this endpoint's return
	go h.UpgradeAuth(context.Background(), user.Password, []byte(authAttempt.Password))

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
		GetLog(r).Error(err)
		return
	}

	GetLog(r).WithField("user", user.ID).Info("Password login success")

}

// UpgradeAuth checks if the password's bcrypt cost is up to current standards, and rehashes it if it's not,
// overwriting the existing creds with a new higher cost version.
func (h *PasswordAuthHandler) UpgradeAuth(ctx context.Context, pass *Password, plaintext []byte) {
	// setup logs
	log := logrus.NewEntry(logrus.StandardLogger()).WithFields(map[string]interface{}{
		"component":   "auth_upgrade",
		"user":        pass.UserID,
		"target_cost": h.PasswordCost,
	})

	cost, err := bcrypt.Cost(pass.Hash)
	if err != nil {
		log.Error(err)
		return
	}

	if cost >= h.PasswordCost {
		// Nothing to be done
		return
	}

	newCreds, err := bcrypt.GenerateFromPassword(plaintext, h.PasswordCost)
	if err != nil {
		log.Error(err)
		return
	}

	// new user object to hold the new creds
	updatedPass := &Password{
		UserID: pass.UserID,
		Hash:   newCreds,
	}

	// Store new creds
	err = h.DB.SetPassword(ctx, updatedPass)
	if err != nil {
		log.WithField("component", "database").Error(err)
	}

	log.Info("Password hash upgraded")
}
