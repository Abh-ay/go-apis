package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	//connect to a PostgreSQL database
	// Replace the connection details (user, dbname, password, host) with your own
	db, err := sqlx.Connect("postgres", "user=postgres dbname=GoDB sslmode=disable password=Qweqwe@1111 host=localhost")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

}
