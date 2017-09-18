package main

// PasswordLogin is the client model for a password auth attempt
type PasswordLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordAuthReply struct {
	Token string `json:"token"`
}

// User is the database model for an auth user
type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Credentials []byte `json:"-"` // never allow this to be rendered to JSON
}
