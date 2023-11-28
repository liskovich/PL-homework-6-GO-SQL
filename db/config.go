package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/api/model"
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

	if isDBEmpty(db) {
		createTables(db)
		seedData(db)
	}
	return db
}

func isDBEmpty(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
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
	// add some default user
	db.Exec("INSERT INTO users (name, email, password) VALUES ('Admin', 'admin@example.com', 'password_hash');")
	row := db.QueryRow(SelectUserByEmailQuery, "admin@example.com")
	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		log.Fatal(err)
	} else {
		// insert some beers
		b1 := fmt.Sprintf(
			"INSERT INTO beers (name, description, thumbnail, author_id) VALUES ('Heineken', 'There is more behind the star', 'https://img.merkandi.com/imgcache/420x320/images/offer/2019/12/05/heineken-lager-1575545711-1575545939.jpeg', %v);",
			user.ID,
		)
		b2 := fmt.Sprintf(
			"INSERT INTO beers (name, description, thumbnail, author_id) VALUES ('Carlsberg', 'Probably the best beer in the world', 'https://auziliquor.com.au/cdn/shop/products/Carlsberg-elephant-premium-strong-beer_900x.jpg?v=1650449387', %v);",
			user.ID,
		)
		b3 := fmt.Sprintf(
			"INSERT INTO beers (name, description, thumbnail, author_id) VALUES ('Cesu PREMIUM original', 'Distinguished by its soft golden color and sparkling refreshment', 'https://www.cesualus.lv/wp-content/uploads/2016/02/CesuPremiumOriginalaisBotPel.jpg', %v);",
			user.ID,
		)
		db.Exec(b1)
		db.Exec(b2)
		db.Exec(b3)
	}
}
