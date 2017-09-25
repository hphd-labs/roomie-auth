package main

type Config struct {
	Port         string `default:"http"`
	PasswordCost int    `default:"10"`
	DatabaseUrl  string `split_words:"true" default:"postgres://postgres:root@localhost:5432/postgres?sslmode=disable"`
}
