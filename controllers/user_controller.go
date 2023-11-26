package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/api/model"
	"example.com/api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController struct {
	authService service.AuthService
}

func NewUserController(service service.AuthService) *UserController {
	return &UserController{
		authService: service,
	}
}

func (ctrl *UserController) RegisterHandler(ctx *gin.Context) {
	body := model.UserMutate{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.authService.Register(body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{})
}

func (ctrl *UserController) LoginHandler(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}
	user := ctrl.authService.Login(body.Email, body.Password)

	// generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// save to cookies
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenStr, 3600*24*30, "", "", false, true)

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{})
}

// TODO: remove in prod
func (ctrl *UserController) ValidateHandler(ctx *gin.Context) {
	// get user from session
	user, _ := ctx.Get("user")

	// access concrete fields
	fmt.Println(user.(model.User).Email)

	ctx.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
