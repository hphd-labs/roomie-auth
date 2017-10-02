package main

type Config struct {
	Port         string `default:"http"`
	PasswordCost int    `default:"12" split_words:"true"`
	DatabaseUrl  string `default:"postgres://postgres:root@localhost:5432/postgres?sslmode=disable" split_words:"true"`
	ShutdownTime int    `default:"30" split_words:"true"`
	WebappOrigin string `default:"roomie-webapp.herokuapp.com" split_words:"true"`
}
