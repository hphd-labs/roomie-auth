package main

import (
	"context"
	"errors"
)

var (
	errUserNotFound = errors.New("User not found")
)

// password is 'password'
// TODO delete this
var dummyUser = User{
	ID:          "1",
	Username:    "admin",
	Credentials: []byte("$2a$10$8rCSa.Tk5L9Qyp4aG4dyP.tVkkadD9B1lV5Pe98QSPkz9lFw6paOa"),
}

type AuthDatabase struct {
}

func (db *AuthDatabase) SetUserCredentials(ctx context.Context, user *User) error {
	if user.ID != dummyUser.ID {
		return errUserNotFound
	}

	dummyUser.Credentials = user.Credentials
	return nil
}

func (db *AuthDatabase) GetUserByName(ctx context.Context, user *User) error {
	if user.Username != dummyUser.Username {
		return errUserNotFound
	}

	user.Credentials = dummyUser.Credentials
	user.ID = dummyUser.ID
	return nil
}
