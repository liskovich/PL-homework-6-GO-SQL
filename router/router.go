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

	// API routes
	// TODO: restructure the user routes to be more API-like
	userRouter.POST(
		"/register",
		userController.RegisterHandler,
	)
	userRouter.POST(
		"/login",
		userController.LoginHandler,
	)
	userRouter.POST(
		"/logout",
		userController.LogoutHandler,
	)
	userRouter.GET(
		"/comments",
		auth,
		userController.UserCommentsHandler,
	)
	userRouter.GET(
		"/beers",
		auth,
		userController.UserBeersHandler,
	)

	beerRouter := baseRouter.Group("/beers")

	beerRouter.POST(
		"/",
		auth,
		beerController.CreateHandler,
	)
	beerRouter.GET(
		"/",
		beerController.GetAllHandler,
	)
	beerRouter.PUT(
		"/:beerId",
		auth,
		beerController.UpdateHandler,
	)
	beerRouter.DELETE(
		"/:beerId",
		auth,
		beerController.DeleteHandler,
	)
	beerRouter.GET(
		"/:beerId",
		optionalUserDetail,
		beerController.GetByIdHandler,
	)

	beerRouter.POST(
		"/:beerId/comment",
		auth,
		beerController.CommentHandler,
	)
	beerRouter.POST(
		"/:beerId/upvote",
		auth,
		beerController.UpvoteHandler,
	)
	beerRouter.POST(
		"/:beerId/downvote",
		auth,
		beerController.DownvoteHandler,
	)

	// UI routes
	// Serve static files (CSS, JavaScript, images, etc.)
	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*.tmpl")
	uiRouter := router.Group("/")

	uiRouter.GET(
		"/",
		optionalUserDetail,
		uiController.Index,
	)
	uiRouter.GET(
		"/dashboard",
		auth,
		uiController.UserDashboard,
	)

	// auth
	uiRouter.GET(
		"/register",
		uiController.RegisterGET,
	)
	uiRouter.GET(
		"/login",
		uiController.LoginGET,
	)
	uiRouter.POST(
		"/register",
		uiController.RegisterPOST,
	)
	uiRouter.POST(
		"/login",
		uiController.LoginPOST,
	)
	uiRouter.POST(
		"/logout",
		uiController.LogoutPOST,
	)

	// basic beer CRUD
	uiRouter.GET(
		"/beers",
		optionalUserDetail,
		uiController.BeersList,
	)
	uiRouter.GET(
		"/beers/create",
		auth,
		uiController.BeersCreateGET,
	)
	uiRouter.POST(
		"/beers/create",
		auth,
		uiController.BeersCreatePOST,
	)
	uiRouter.GET(
		"/beers/:beerId",
		optionalUserDetail,
		uiController.BeersDetail,
	)
	uiRouter.GET(
		"/beers/:beerId/edit",
		auth,
		uiController.BeersEditGET,
	)
	uiRouter.POST(
		"/beers/:beerId/edit",
		auth,
		uiController.BeersEditPOST,
	)
	uiRouter.POST(
		"/beers/:beerId/delete",
		auth,
		uiController.BeersDeletePOST,
	)

	// additional features
	uiRouter.POST(
		"/beers/:beerId/comment",
		auth,
		uiController.CommentPOST,
	)
	uiRouter.POST(
		"/beers/:beerId/upvote",
		auth,
		uiController.UpvotePOST,
	)
	uiRouter.POST(
		"/beers/:beerId/downvote",
		auth,
		uiController.DownvotePOST,
	)

	return router
}
