package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func (ctrl *UIController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index", gin.H{})
}

func (ctrl *UIController) UserDashboard(ctx *gin.Context) {
	// TODO: pass data to template
	ctx.HTML(http.StatusOK, "dashboard", gin.H{})
}

// auth endpoints
func (ctrl *UIController) RegisterGET(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register", gin.H{})
}

func (ctrl *UIController) LoginGET(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login", gin.H{})
}

func (ctrl *UIController) RegisterPOST(ctx *gin.Context) {
	body := model.UserMutate{
		Name:     ctx.PostForm("username"),
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.authService.Register(body)
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) LoginPOST(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	body.Email = ctx.PostForm("email")
	body.Password = ctx.PostForm("password")
	if ctx.Bind(&body) != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
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
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// save to cookies
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenStr, 3600*24*30, "", "", false, true)
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) LogoutPOST(ctx *gin.Context) {
	// clear the Authorization cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", "", -1, "", "", false, true)
	ctx.Redirect(http.StatusFound, "/")
}

// beer basic CRUD
func (ctrl *UIController) BeersList(ctx *gin.Context) {
	beers := ctrl.beerService.FindAll()
	ctx.HTML(http.StatusOK, "beers", gin.H{
		"Beers": beers,
	})
}

func (ctrl *UIController) BeersDetail(ctx *gin.Context) {
	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}
	beer := ctrl.beerService.FindById(uint(beerID))
	if beer == nil {
		ctx.HTML(http.StatusNotFound, "error", gin.H{
			"error": "Your requested beer was not found",
		})
		return
	}

	var hasUpvoted bool = false
	currentUser, usrKeyExists := ctx.Get("user")
	if usrKeyExists && currentUser != nil {
		hasUpvoted = ctrl.upvoteService.CheckIfUserUpvoted(model.Upvote{
			BeerID: beer.ID,
			UserID: currentUser.(*model.User).ID,
		})
	}

	ctx.HTML(http.StatusOK, "beers_detail", gin.H{
		"Beer":        *beer,
		"UserUpvoted": hasUpvoted,
	})
}

func (ctrl *UIController) BeersCreateGET(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "beers_create", gin.H{})
}

func (ctrl *UIController) BeersCreatePOST(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists && currentUser == nil {
		ctx.HTML(http.StatusUnauthorized, "error", gin.H{
			"error": "You must be logged in to create a beer",
		})
		return
	}

	var body model.BeerMutate
	body.Name = ctx.PostForm("name")
	body.Description = ctx.PostForm("description")
	body.Thumbnail = ctx.PostForm("thumbnail")
	body.AuthorId = currentUser.(*model.User).ID
	if ctx.Bind(&body) != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.beerService.Create(body)
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) BeersEditGET(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists && currentUser == nil {
		ctx.HTML(http.StatusUnauthorized, "error", gin.H{
			"error": "You must be logged in to create a beer",
		})
		return
	}

	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}
	beer := ctrl.beerService.FindById(uint(beerID))
	if beer == nil {
		ctx.HTML(http.StatusNotFound, "error", gin.H{
			"error": "Your requested beer was not found",
		})
		return
	}

	if beer.AuthorId != currentUser.(*model.User).ID {
		ctx.HTML(http.StatusUnauthorized, "error", gin.H{
			"error": "Your have to be the author of the beer to edit it",
		})
		return
	}
	ctx.HTML(http.StatusOK, "beers_edit", gin.H{
		"Beer": beer,
	})
}

func (ctrl *UIController) BeersEditPOST(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists && currentUser == nil {
		ctx.HTML(http.StatusUnauthorized, "error", gin.H{
			"error": "You must be logged in to create a beer",
		})
		return
	}

	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}

	var body model.BeerMutate
	body.Name = ctx.PostForm("name")
	body.Description = ctx.PostForm("description")
	body.Thumbnail = ctx.PostForm("thumbnail")
	body.AuthorId = currentUser.(*model.User).ID
	if ctx.Bind(&body) != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.beerService.Update(uint(beerID), body)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/beers/%v", beerID))
}

func (ctrl *UIController) BeersDeletePOST(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists && currentUser == nil {
		ctx.HTML(http.StatusUnauthorized, "error", gin.H{
			"error": "You must be logged in to create a beer",
		})
		return
	}

	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error", gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}

	ctrl.beerService.Delete(
		uint(beerID),
		currentUser.(*model.User).ID,
	)
	ctx.Redirect(http.StatusFound, "/beers")
}

// additional features
func (ctrl *UIController) CommentPOST(ctx *gin.Context) {
	// TODO: redirect to the commented beer
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) UpvoteDownvotePOST(ctx *gin.Context) {
	// TODO: redirect to the upvoted / downvoted beer
	ctx.Redirect(http.StatusFound, "/beers")
}
