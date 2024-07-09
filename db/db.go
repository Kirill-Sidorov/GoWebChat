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
		log.Fatal("environment variable db_username not found in .env")
	}
	dbPassword, found := os.LookupEnv("db_password")
	if !found {
		log.Fatal("environment variable db_password not found in .env")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@postgres:5432/webchatdb?sslmode=disable", dbUser, dbPassword)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("DB connection close")
	}
}
