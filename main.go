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

	// init repositories
	userRepository := db.NewUserRepository(database)
	beerRepository := db.NewBeerRepository(database)
	commentRepository := db.NewCommentRepository(database)
	upvoteRepository := db.NewUpvoteRepository(database)

	// init service layer
	beerService := service.NewBeerService(
		beerRepository,
		commentRepository,
		validate,
	)
	authService := service.NewAuthService(
		userRepository,
		validate,
	)
	commentService := service.NewCommentService(
		commentRepository,
		beerRepository,
		validate,
	)
	upvoteService := service.NewUpvoteService(
		upvoteRepository,
		validate,
	)

	// init controllers
	userController := controllers.NewUserController(
		authService,
		commentService,
		beerService,
	)
	beerController := controllers.NewBeerController(
		authService,
		beerService,
		commentService,
		upvoteService,
	)
	uiController := controllers.NewUIController(
		authService,
		beerService,
		commentService,
		upvoteService,
	)

	// init server
	router := router.NewRouter(
		userController,
		beerController,
		uiController,
		&authService,
	)

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
