package service

import (
	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type BeerService interface {
	Create(beer model.BeerMutate)
	Update(beerID uint, beer model.BeerMutate)
	Delete(beerID uint)
	FindById(beerID uint) model.BeerDetailed
	FindAll() []*model.BeerCompact

	Upvote(beerID uint, userID uint)
	Downvote(beerID uint, userID uint)
	Comment(comment model.CommentMutate)
}

type BeerServiceImpl struct {
	BeerRepository    db.BeerRepository
	CommentRepository db.CommentRepository
	UpvoteRepository  db.UpvoteRepository
	Validate          *validator.Validate
}

func NewBeerService(
	beerRepository db.BeerRepository,
	commentRepository db.CommentRepository,
	upvoteRepository db.UpvoteRepository,
	validate *validator.Validate) BeerService {
	return &BeerServiceImpl{
		BeerRepository:    beerRepository,
		CommentRepository: commentRepository,
		UpvoteRepository:  upvoteRepository,
		Validate:          validate,
	}
}

func (brService *BeerServiceImpl) Create(beer model.BeerMutate) {
	err := brService.Validate.Struct(beer)
	if err != nil {
		panic(err)
	}
	brService.BeerRepository.CreateBeer(beer)
}

func (brService *BeerServiceImpl) Delete(beerID uint) {
	err := brService.BeerRepository.DeleteBeer(beerID)
	if err != nil {
		panic(err)
	}
}

func (brService *BeerServiceImpl) FindAll() []*model.BeerCompact {
	result, err := brService.BeerRepository.GetAllBeers()
	if err != nil {
		panic(err)
	}
	return result
}

func (brService *BeerServiceImpl) FindById(beerID uint) model.BeerDetailed {
	result, err := brService.BeerRepository.GetBeerById(beerID)
	if err != nil {
		panic(err)
	}
	comments, cmntErr := brService.CommentRepository.GetAllBeerComments(beerID)
	if cmntErr != nil {
		panic(cmntErr)
	}
	beerResponse := model.BeerDetailed{
		ID:           result.ID,
		Name:         result.Name,
		Description:  result.Description,
		Thumbnail:    result.Thumbnail,
		CommentCount: result.CommentCount,
		UpvoteCount:  result.UpvoteCount,
		Comments:     comments,
	}
	return beerResponse
}

func (brService *BeerServiceImpl) Update(beerID uint, beer model.BeerMutate) {
	beerToUpdate, err := brService.BeerRepository.GetBeerById(beerID)
	if err != nil || beerToUpdate == nil {
		panic(err)
	}
	brService.BeerRepository.UpdateBeer(beerID, beer)
}

func (brService *BeerServiceImpl) Comment(comment model.CommentMutate) {
	// TODO: check if user is authenticated and beer exists
	err := brService.CommentRepository.CreateComment(comment)
	if err != nil {
		panic(err)
	}
}

func (brService *BeerServiceImpl) Upvote(beerID uint, userID uint) {
	upvote := model.Upvote{
		UserID: userID,
		BeerID: beerID,
	}
	exists, err := brService.UpvoteRepository.CheckUpvoteExists(upvote)
	if exists || err != nil {
		panic(err)
	}
	upvoteErr := brService.UpvoteRepository.CreateUpvote(upvote)
	if upvoteErr != nil {
		panic(upvoteErr)
	}
}

func (brService *BeerServiceImpl) Downvote(beerID uint, userID uint) {
	downvote := model.Upvote{
		UserID: userID,
		BeerID: beerID,
	}
	exists, err := brService.UpvoteRepository.CheckUpvoteExists(downvote)
	if exists || err != nil {
		panic(err)
	}
	downvoteErr := brService.UpvoteRepository.DeleteUpvote(downvote)
	if downvoteErr != nil {
		panic(downvoteErr)
	}
}
