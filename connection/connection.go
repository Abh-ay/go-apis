package dbConnection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	connStr := "postgres://postgres:Qweqwe@1111@localhost/GoDB?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("The database is connected")

}
