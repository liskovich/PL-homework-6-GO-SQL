package main

import (
	"net/http"
	"time"

	"example.com/api/db"
	"example.com/api/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	database := db.ConnectDB()
	validate := validator.New()

	userRepository := db.NewUserRepository(database)
	beerRepository := db.NewBeerRepository(database)
	commentRepository := db.NewCommentRepository(database)
	upvoteRepository := db.NewUpvoteRepository(database)

	beerService := service.NewBeerService(beerRepository, commentRepository, upvoteRepository, validate)
	authService := service.NewAuthService(userRepository, validate)

	router := gin.Default()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello, world")
	})

	server := &http.Server{
		Addr:           ":8888",
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
