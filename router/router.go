package router

import (
	"example.com/api/controllers"
	"example.com/api/middleware"
	"example.com/api/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controllers.UserController, authService *service.AuthService) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.HandlePanic())

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})
	baseRouter := router.Group("/api")
	userRouter := baseRouter.Group("/user")

	userRouter.POST("/register", userController.RegisterHandler)
	userRouter.POST("/login", userController.LoginHandler)
	userRouter.GET("/validate", middleware.AuthMiddleware(*authService), userController.ValidateHandler)
	return router
}
