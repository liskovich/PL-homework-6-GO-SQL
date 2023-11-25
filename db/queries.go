package db

const (
	CreateUserTableQuery     = "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT UNIQUE NOT NULL, password TEXT NOT NULL);"
	CreateBeerTableQuery     = "CREATE TABLE IF NOT EXISTS beers (id SERIAL PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL, thumbnail TEXT NOT NULL);"
	CreateCommentsTableQuery = "CREATE TABLE IF NOT EXISTS comments (id SERIAL PRIMARY KEY, author_id INTEGER REFERENCES users(id), content TEXT NOT NULL, created_date DATE NOT NULL, beer_id INTEGER REFERENCES beers(id));"
	CreateUpvotesTableQuery  = "CREATE TABLE IF NOT EXISTS upvotes (id SERIAL PRIMARY KEY, user_id INTEGER REFERENCES users(id), beer_id INTEGER REFERENCES beers(id), UNIQUE (user_id, beer_id));"

	InsertUserQuery     = "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	SelectUserByIdQuery = "SELECT id, name, email, password FROM users WHERE id = $1"
	UpdateUserQuery     = "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"

	InsertBeerQuery     = "INSERT INTO beers (name, description, thumbnail) VALUES ($1, $2, $3)"
	SelectBeerByIdQuery = "SELECT id, name, description, thumbnail FROM beers WHERE id = $1" // TODO: join with comments and upvotes
	SelectAllBeersQuery = "SELECT id, name, description, thumbnail FROM beers"               // TODO: join with comments and upvotes
	UpdateBeerQuery     = "UPDATE beers SET name = $1, description = $2, thumbnail = $3 WHERE id = $4"

	InsertCommentQuery        = "INSERT INTO comments (author_id, content, created_date, beer_id) VALUES ($1, $2, $3, $4)"
	SelectAllCommentsByUserId = "SELECT id, author_id, content, created_date, beer_id FROM comments WHERE author_id = $1"

	InsertUpvoteQuery        = "INSERT INTO upvotes (user_id, beer_id) VALUES ($1, $2)"
	DeleteUpvoteQuery        = "DELETE FROM upvotes WHERE user_id = $1 AND beer_id = $2"
	CheckIfUpvoteExistsQuery = "SELECT COUNT(*) FROM upvotes WHERE user_id = $1 AND beer_id = $2"
)
