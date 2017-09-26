package main

import (
	"context"
	"github.com/pkg/errors"
)

type MockDatabase struct {
	users     []*User
	passwords []*Password
	calls     map[string]int
}

func newMockDatabase() *MockDatabase {
	m := &MockDatabase{
		users:     make([]*User, 0, 1),
		passwords: make([]*Password, 0, 1),
		calls:     make(map[string]int),
	}

	dummyUser := &User{
		ID:       "1",
		Username: "test",
		Password: &Password{
			UserID: "1",
			Hash:   []byte("$2a$04$bERAY/5FeoucwWTSGWvAwuwgzWYsVTM1hYtKApmg74qxxzxGG/QQ6"),
		},
	}

	m.users = append(m.users, dummyUser)
	m.passwords = append(m.passwords, dummyUser.Password)

	return m
}

func (m *MockDatabase) SetPassword(_ context.Context, target *Password) error {
	m.calls["SetPassword"] = m.calls["SetPassword"] + 1
	for _, pass := range m.passwords {
		if pass.UserID == target.UserID {
			pass.Hash = target.Hash
			return nil
		}
	}

	return errors.New("Password not found")
}

func (m *MockDatabase) GetUserByName(_ context.Context, target *User) error {
	m.calls["GetUserByName"] = m.calls["GetUserByName"] + 1
	for _, user := range m.users {
		if user.Username == target.Username {
			*target = *user
			return nil
		}
	}

	return ErrUserNotFound
}
