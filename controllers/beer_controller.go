package controllers

import (
	"net/http"
	"strconv"

	"example.com/api/model"
	"example.com/api/service"
	"github.com/gin-gonic/gin"
)

type BeerController struct {
	authService    service.AuthService
	beerService    service.BeerService
	commentService service.CommentService
	upvoteService  service.UpvoteService
}

func NewBeerController(
	authSrvc service.AuthService,
	beerSrvc service.BeerService,
	cmntSrvc service.CommentService,
	upvtSrvc service.UpvoteService,
) *BeerController {
	return &BeerController{
		authService:    authSrvc,
		beerService:    beerSrvc,
		commentService: cmntSrvc,
		upvoteService:  upvtSrvc,
	}
}

func (ctrl *BeerController) CreateHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to create a beer",
		})
		return
	}
	var body model.BeerMutate
	body.AuthorId = currentUser.(model.User).ID
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	ctrl.beerService.Create(body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, gin.H{})
}

func (ctrl *BeerController) UpdateHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to update a beer",
		})
		return
	}
	var body model.BeerMutate
	body.AuthorId = currentUser.(model.User).ID
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}
	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}

	ctrl.beerService.Update(uint(beerID), body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{})
}

func (ctrl *BeerController) DeleteHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to delete a beer",
		})
		return
	}
	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}
	ctrl.beerService.Delete(uint(beerID), currentUser.(model.User).ID)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (ctrl *BeerController) GetAllHandler(ctx *gin.Context) {
	beers := ctrl.beerService.FindAll()
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"data": beers,
	})
}

func (ctrl *BeerController) GetByIdHandler(ctx *gin.Context) {
	beerIDStr := ctx.Param("beerId")
	beerID, err := strconv.Atoi(beerIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid beer ID provided",
		})
		return
	}
	beer := ctrl.beerService.FindById(uint(beerID))
	if beer == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var hasUpvoted bool = false
	currentUser, usrKeyExists := ctx.Get("user")
	if usrKeyExists && currentUser != nil {
		hasUpvoted = ctrl.upvoteService.CheckIfUserUpvoted(model.Upvote{
			BeerID: beer.ID,
			UserID: currentUser.(model.User).ID,
		})
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"data":         beer,
		"user_upvoted": hasUpvoted,
	})
}

func (ctrl *BeerController) CommentHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to comment",
		})
		return
	}
	var body model.CommentMutate
	body.AuthorID = currentUser.(model.User).ID
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	beerToComment := ctrl.beerService.FindById(body.BeerID)
	if beerToComment == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "The beer you wanted to comment was not found",
		})
		return
	}
	ctrl.commentService.Create(body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, gin.H{})
}

func (ctrl *BeerController) UpvoteHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to upvote",
		})
		return
	}
	var body model.Upvote
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	beerToUpvote := ctrl.beerService.FindById(body.BeerID)
	if beerToUpvote == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "The beer you wanted to upvote was not found",
		})
		return
	}
	ctrl.upvoteService.Upvote(body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, gin.H{})
}

func (ctrl *BeerController) DownvoteHandler(ctx *gin.Context) {
	currentUser, usrKeyExists := ctx.Get("user")
	if !usrKeyExists || currentUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to downvote",
		})
		return
	}
	var body model.Upvote
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse request body",
		})
		return
	}

	beerToDownvote := ctrl.beerService.FindById(body.BeerID)
	if beerToDownvote == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "The beer you wanted to downvote was not found",
		})
		return
	}
	ctrl.upvoteService.Downvote(body)
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusCreated, gin.H{})
}
