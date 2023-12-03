# Go CRUD with Postgres (SQL) database

I built **Beer Hunt** - it is like [Product Hunt](https://www.producthunt.com/), but for beers.

## Implemented functions:

- **Web UI** - this is a web app, built using Go Gin web framework.
- **Auth** - login and register using JWT token, password hashing.
- **Beer list** - displays a list of beers.
- **Beer comments** - users can post comments on beers (must be logged in).
- **Beer upvote / downvote** - users can upvote and downvote beers (must be logged in).
- **Dashboard** - where a logged in user can see their created beers and posted comments.
- **Create, edit, delete beers** - logged in users have access to full beer CRUD functionality (only if they are the authors of this beer).

## How to run?

**IMPORTANT NOTE:**
This program was developed in VS code.

**Prerequisites:**

- You need to have [go compiler](https://go.dev/doc/install) installed and set up on your machine.
- Docker desktop needs to be installed to run the Postgres database in a container.
- Docker compose needs to be installed.

**Running program:**

1. Before running the program, create a copy of file `.env.example` and rename it to `.env`.
2. It should have the following contents:
```env
PORT=":8888"
SECRET=<YOUR_SECRET>

DB_USER=<YOUR_USER>
DB_PASSWORD=<YOUR_PASSWORD>
DB_NAME=beerhunt
DB_HOST=localhost
```
3. Replace the `<YOUR_USER>` and `<YOUR_PASSWORD>` with your desired database credentials. These would be picked up when building the Docker container. Also, replace `<YOUR_SECRET>` with some alphanumeric character array (this will be used for the JWT auth).
4. Go to the root of the project and run `docker compose up --build` to and wait for it to start the Postgres container.
5. Run `go run .` to build and run the go app.
6. Navigate to [http://127.0.0.1:8888](http://127.0.0.1:8888). 
7. If everything is OK, the index page of application should be visible.

## Trying out the app

To test the app, the workflow would be following:

1. Click `Register` and create a new user.
2. Click `Log In` and authenticate your newly created user.
3. You can now use the full beer CRUD functionality.
4. Click `Log Out` when finished work.