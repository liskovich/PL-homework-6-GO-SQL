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
	// TODO: pass data to template
	ctx.HTML(http.StatusOK, "beers_detail.tmpl", gin.H{})
}

func (ctrl *UIController) BeersCreateGET(ctx *gin.Context) {
	// TODO: pass data to template
	ctx.HTML(http.StatusOK, "beers_create.tmpl", gin.H{})
}

func (ctrl *UIController) BeersCreatePOST(ctx *gin.Context) {
	// TODO: redirect to newly created beer
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) BeersEditGET(ctx *gin.Context) {
	// TODO: pass data to template
	ctx.HTML(http.StatusOK, "beers_edit.tmpl", gin.H{})
}

func (ctrl *UIController) BeersEditPOST(ctx *gin.Context) {
	// TODO: redirect to updated beer
	ctx.Redirect(http.StatusFound, "/beers")
}

func (ctrl *UIController) BeersDeletePOST(ctx *gin.Context) {
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
