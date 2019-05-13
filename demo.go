package main

import (
	"database/sql"
	"fmt"
	_ "net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "oluwatobiloba"
	password = "nat1234"
	dbname   = "db_demo"
)

func help() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	sqlStatement := `INSERT INTO users (age, email, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`

	id := 0
	err = db.QueryRow(sqlStatement, 33, "jkelvin@example.com", "jamie", "kelvin").Scan(&id)

	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}
