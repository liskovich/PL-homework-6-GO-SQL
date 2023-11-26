package db

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	connStr := "postgres://username:password@localhost/dbname?sslmode=disable"
	connStr = strings.Replace(connStr, "username", username, 1)
	connStr = strings.Replace(connStr, "password", password, 1)
	connStr = strings.Replace(connStr, "localhost", dbHost, 1)
	connStr = strings.Replace(connStr, "dbname", dbName, 1)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables(db)
	seedData(db)
	return db
}

func createTables(db *sql.DB) {
	for _, query := range []string{CreateUserTableQuery, CreateBeerTableQuery, CreateCommentsTableQuery, CreateUpvotesTableQuery} {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedData(db *sql.DB) {
	// TODO: add seed data if needed
	// Example: db.Exec("INSERT INTO users (name, email, password) VALUES ('Alice', 'alice@example.com', 'password')")
}
