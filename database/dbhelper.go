package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // importing the postgres helper

	"github.com/joho/godotenv"
)

var (
	db       *sql.DB
	host     string
	port     string
	user     string
	password string
	dbname   string
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file", err.Error(), err)
	}
}

// GetDb :
func GetDb() *sql.DB {
	host := os.Getenv("db_host")
	port := os.Getenv("db_port")
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err.Error(), err)
	}

	if err = db.Ping(); err != nil {
		log.Println(err.Error(), err)
	}

	fmt.Println("Database connnected!")
	return db
}
