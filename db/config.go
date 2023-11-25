package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// Extract credentials to env file
	username := "postgres"
	password := "postgres"
	connStr := "postgres://username:password@localhost/dbname?sslmode=disable"
	connStr = strings.Replace(connStr, "username", username, 1)
	connStr = strings.Replace(connStr, "password", password, 1)

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
