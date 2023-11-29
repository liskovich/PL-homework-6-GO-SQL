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

type UIController struct {
	authService    service.AuthService
	beerService    service.BeerService
	commentService service.CommentService
	upvoteService  service.UpvoteService
}

func NewUIController(
	authSrvc service.AuthService,
	beerSrvc service.BeerService,
	cmntSrvc service.CommentService,
	upvtSrvc service.UpvoteService,
) *UIController {
	return &UIController{
		authService:    authSrvc,
		beerService:    beerSrvc,
		commentService: cmntSrvc,
		upvoteService:  upvtSrvc,
	}
}

// auth endpoints
func (ctrl *UIController) RegisterGET(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.tmpl", gin.H{})
}

func (ctrl *UIController) LoginGET(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func (ctrl *UIController) RegisterPOST(ctx *gin.Context) {
	body := model.UserMutate{
		Name:     ctx.PostForm("username"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.authService.Register(body)
	ctx.Redirect(http.StatusFound, "/")
}

func (ctrl *UIController) LoginPOST(ctx *gin.Context) {
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

// beer basic CRUD
func (ctrl *UIController) Index(ctx *gin.Context) {
	// TODO: pass data to template
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// additional features
