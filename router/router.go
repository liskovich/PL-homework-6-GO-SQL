package router

import (
	"example.com/api/controllers"
	"example.com/api/middleware"
	"example.com/api/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	userController *controllers.UserController,
	beerController *controllers.BeerController,
	uiController *controllers.UIController,
	authService *service.AuthService,
) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.HandlePanic())

	baseRouter := router.Group("/api")
	userRouter := baseRouter.Group("/user")

	auth := middleware.AuthMiddleware(*authService)
	optionalUserDetail := middleware.UserDetailMiddleware(*authService)

	// TODO: restructure the user routes to be more API-like
	userRouter.POST("/register", userController.RegisterHandler)
	userRouter.POST("/login", userController.LoginHandler)
	// TODO: add log out
	userRouter.GET("/comments", auth, userController.UserCommentsHandler)
	userRouter.GET("/beers", auth, userController.UserBeersHandler)

	beerRouter := baseRouter.Group("/beers")
	beerRouter.POST("/", auth, beerController.CreateHandler)
	beerRouter.GET("/", beerController.GetAllHandler)
	beerRouter.PUT("/:beerId", auth, beerController.UpdateHandler)
	beerRouter.DELETE("/:beerId", auth, beerController.DeleteHandler)
	beerRouter.GET("/:beerId", optionalUserDetail, beerController.GetByIdHandler)

	beerRouter.POST("/:beerId/comment", auth, beerController.CommentHandler)
	beerRouter.POST("/:beerId/upvote", auth, beerController.UpvoteHandler)
	beerRouter.POST("/:beerId/downvote", auth, beerController.DownvoteHandler)

	// UI routes
	router.LoadHTMLGlob("templates/*.tmpl")
	uiRouter := router.Group("/")

	uiRouter.GET("/", uiController.Index)

	uiRouter.GET("/register", uiController.RegisterGET)
	uiRouter.GET("/login", uiController.LoginGET)
	uiRouter.POST("/register", uiController.RegisterPOST)
	uiRouter.POST("/login", uiController.LoginPOST)

	return router
}
