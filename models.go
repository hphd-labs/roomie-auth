package main

// PasswordLogin is the client model for a password auth attempt
type PasswordLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
