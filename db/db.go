package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() {
	dbUser, found := os.LookupEnv("db_username")
	if !found {
		log.Panic("environment variable db_username not found in .env")
	}

	dbPassword, found := os.LookupEnv("db_password")
	if !found {
		log.Panic("environment variable db_password not found in .env")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=webchat sslmode=disable", dbUser, dbPassword)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
}