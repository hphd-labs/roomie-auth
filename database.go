package main

import (
	"context"
	"errors"
	"github.com/go-pg/pg"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

type AuthDatabase struct {
	Database *pg.DB
}

func (db *AuthDatabase) SetPassword(ctx context.Context, pass *Password) error {
	_, err := db.Database.
		WithContext(ctx).
		Model(pass).
		Set("hash = ?hash").
		Update()

	return err
}

func (db *AuthDatabase) GetUserByName(ctx context.Context, user *User) error {
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
