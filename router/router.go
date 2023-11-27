package router

import (
	"example.com/api/controllers"
	"example.com/api/middleware"
	"example.com/api/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userController *controllers.UserController,
	beerController *controllers.BeerController,
	authService *service.AuthService,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.HandlePanic())

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})
	baseRouter := router.Group("/api")
	userRouter := baseRouter.Group("/user")

	// TODO: restructure the user routes to be more API-like
	userRouter.POST("/register", userController.RegisterHandler)
	userRouter.POST("/login", userController.LoginHandler)
	userRouter.GET("/comments", middleware.AuthMiddleware(*authService), userController.UserCommentsHandler)
	userRouter.GET("/beers", middleware.AuthMiddleware(*authService), userController.UserBeersHandler)

	beerRouter := baseRouter.Group("/beers")
	beerRouter.POST("/", middleware.AuthMiddleware(*authService), beerController.CreateHandler)
	beerRouter.POST("/", beerController.GetAllHandler)
	beerRouter.PUT("/:beerId", middleware.AuthMiddleware(*authService), beerController.UpdateHandler)
	beerRouter.DELETE("/:beerId", middleware.AuthMiddleware(*authService), beerController.DeleteHandler)
	beerRouter.GET("/:beerId", beerController.GetByIdHandler)

	return router
}
