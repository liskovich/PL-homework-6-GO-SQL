package db

const (
	CreateUserTableQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY, 
			name TEXT UNIQUE NOT NULL, 
			email TEXT UNIQUE NOT NULL, 
			password TEXT NOT NULL
		);`
	CreateBeerTableQuery = `
		CREATE TABLE IF NOT EXISTS beers (
			id SERIAL PRIMARY KEY, 
			name TEXT NOT NULL, 
			description TEXT NOT NULL, 
			thumbnail TEXT NOT NULL,
			author_id INTEGER REFERENCES users(id), 
		);`
	CreateCommentsTableQuery = `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY, 
			author_id INTEGER REFERENCES users(id), 
			content TEXT NOT NULL, 
			created_date DATE NOT NULL, 
			beer_id INTEGER REFERENCES beers(id) ON DELETE CASCADE
		);`
	CreateUpvotesTableQuery = `
		CREATE TABLE IF NOT EXISTS upvotes (
			id SERIAL PRIMARY KEY, 
			user_id INTEGER REFERENCES users(id), 
			beer_id INTEGER REFERENCES beers(id) ON DELETE CASCADE, 
			UNIQUE (user_id, beer_id)
		);`

	InsertUserQuery        = "INSERT INTO users (name, email, password) VALUES ($1, $2, $3);"
	SelectUserByIdQuery    = "SELECT id, name, email, password FROM users WHERE id = $1;"
	SelectUserByEmailQuery = "SELECT id, name, email, password FROM users WHERE email = $1;"
	UpdateUserQuery        = "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4;"

	InsertBeerQuery     = "INSERT INTO beers (name, description, thumbnail, author_id) VALUES ($1, $2, $3, $4);"
	SelectBeerByIdQuery = `
		SELECT 
			b.id,
    	b.name,
    	b.description,
    	b.thumbnail,
			b.author_id,
    	COALESCE(comment_count, 0) AS comment_count,
    	COALESCE(upvote_count, 0) AS upvote_count
		FROM beers b
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY beer_id
		) c ON b.id = c.beer_id
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS upvote_count
			FROM upvotes
			GROUP BY beer_id
		) u ON b.id = u.beer_id
		WHERE b.id = $1;`
	SelectBeersByUserQuery = `
		SELECT 
			b.id,
    	b.name,
    	b.description,
    	b.thumbnail,
			b.author_id,
    	COALESCE(comment_count, 0) AS comment_count,
    	COALESCE(upvote_count, 0) AS upvote_count
		FROM beers b
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY beer_id
		) c ON b.id = c.beer_id
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS upvote_count
			FROM upvotes
			GROUP BY beer_id
		) u ON b.id = u.beer_id
		WHERE b.author_id = $1;`
	SelectAllBeersQuery = `
		SELECT 
			b.id,
    	b.name,
    	b.description,
    	b.thumbnail,
			b.author_id,
    	COALESCE(comment_count, 0) AS comment_count,
    	COALESCE(upvote_count, 0) AS upvote_count
		FROM beers b
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY beer_id
		) c ON b.id = c.beer_id
		LEFT JOIN (
			SELECT beer_id, COUNT(*) AS upvote_count
			FROM upvotes
			GROUP BY beer_id
		) u ON b.id = u.beer_id;`
	UpdateBeerQuery = "UPDATE beers SET name = $1, description = $2, thumbnail = $3 WHERE id = $4;"
	DeleteBeerQuery = "DELETE FROM beers WHERE id = $1;"

	InsertCommentQuery             = "INSERT INTO comments (author_id, content, created_date, beer_id) VALUES ($1, $2, $3, $4);"
	SelectAllCommentsByUserIdQuery = "SELECT id, author_id, content, created_date, beer_id FROM comments WHERE author_id = $1;"
	SelectAllCommentsByBeerIdQuery = "SELECT id, author_id, content, created_date, beer_id FROM comments WHERE beer_id = $1;"

	InsertUpvoteQuery        = "INSERT INTO upvotes (user_id, beer_id) VALUES ($1, $2);"
	DeleteUpvoteQuery        = "DELETE FROM upvotes WHERE user_id = $1 AND beer_id = $2;"
	CheckIfUpvoteExistsQuery = "SELECT COUNT(*) FROM upvotes WHERE user_id = $1 AND beer_id = $2;"
)
