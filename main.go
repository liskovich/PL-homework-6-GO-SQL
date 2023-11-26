package main

import (
	"net/http"
	"os"
	"time"

	"example.com/api/config"
	"example.com/api/controllers"
	"example.com/api/db"
	"example.com/api/router"
	"example.com/api/service"
	"github.com/go-playground/validator/v10"
)

// https://youtu.be/ma7rUS_vW9M?feature=shared

func main() {
	config.LoadEnvVariables()

	database := db.ConnectDB()
	validate := validator.New()

	userRepository := db.NewUserRepository(database)
	beerRepository := db.NewBeerRepository(database)
	commentRepository := db.NewCommentRepository(database)
	upvoteRepository := db.NewUpvoteRepository(database)

	beerService := service.NewBeerService(beerRepository, commentRepository, upvoteRepository, validate)
	authService := service.NewAuthService(userRepository, validate)

	userController := controllers.NewUserController(authService)

	router := router.NewRouter(userController, &authService)

	// TODO: extract to env file
	server := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
