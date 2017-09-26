package main

import (
	"context"
	"errors"
	"github.com/go-pg/pg"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

// AuthDatabase is an interface describing all the functionality an auth
// database backed must implement. This is an interface to allow mock testing
type AuthDatabase interface {
	SetPassword(context.Context, *Password) error
	GetUserByName(context.Context, *User) error
}

type PGAuthDatabase struct {
	Database *pg.DB
}

func (db *PGAuthDatabase) SetPassword(ctx context.Context, pass *Password) error {
	_, err := db.Database.
		WithContext(ctx).
		Model(pass).
		Set("hash = ?hash").
		Update()

	return err
}

func (db *PGAuthDatabase) GetUserByName(ctx context.Context, user *User) error {
	err := db.Database.
		WithContext(ctx).
		Model(user).
		Column("Password").
		Where("username = ?username").
		Select()

	if err == pg.ErrNoRows {
		return ErrUserNotFound
	}
	return err
}
