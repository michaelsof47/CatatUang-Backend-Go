package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDatabase() *sql.DB {
	conn := "host=localhost port=5432 user=postgres password=1234 dbname=CatatUang_DB sslmode=disable"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database Terhubung")
	return db
}
