package main

import (
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var postgresCredentials string
var db *sqlx.DB

func init() {
	b, err := ioutil.ReadFile("config/postgres_credentials.txt")
	panicIfError(err)
	postgresCredentials = string(b)

	cutIndex := len(postgresCredentials) - 1
	if postgresCredentials[cutIndex] == '\n' {
		postgresCredentials = postgresCredentials[:cutIndex]
	}

	db, err = sqlx.Connect("postgres", postgresCredentials)
	panicIfError(err)
}
