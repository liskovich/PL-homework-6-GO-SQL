package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/api/model"
	"example.com/api/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController struct {
	authService    service.AuthService
	commentService service.CommentService
	beerService    service.BeerService
}

func NewUserController(authSrvc service.AuthService, cmntSrvc service.CommentService, beerSrvc service.BeerService) *UserController {
	return &UserController{
		authService:    authSrvc,
		commentService: cmntSrvc,
		beerService:    beerSrvc,
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

func (ctrl *UserController) UserCommentsHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in",
		})
		return
	}

	comments := ctrl.commentService.FindAllUsersComments(currentUser.(model.User).ID)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

func (ctrl *UserController) UserBeersHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged",
		})
		return
	}

	beers := ctrl.beerService.FindByUser(currentUser.(model.User).ID)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"data": beers,
	})
}
