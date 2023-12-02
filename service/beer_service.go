package service

import (
	"database/sql"

	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type BeerService interface {
	Create(beer model.BeerMutate)
	Update(beerID uint, beer model.BeerMutate)
	Delete(beerID uint, userID uint)
	FindById(beerID uint) *model.BeerDetailed
	FindByUser(userID uint) []model.BeerCompact
	FindAll() []model.BeerCompact
}

type BeerServiceImpl struct {
	BeerRepository    db.BeerRepository
	CommentRepository db.CommentRepository
	Validate          *validator.Validate
}

func NewBeerService(
	beerRepository db.BeerRepository,
	commentRepository db.CommentRepository,
	validate *validator.Validate) BeerService {
	return &BeerServiceImpl{
		BeerRepository:    beerRepository,
		CommentRepository: commentRepository,
		Validate:          validate,
	}
}

func (brService *BeerServiceImpl) Create(beer model.BeerMutate) {
	// TODO: determine what does it do
	err := brService.Validate.Struct(beer)
	if err != nil {
		panic(err)
	}
	createErr := brService.BeerRepository.CreateBeer(beer)
	if createErr != nil {
		panic(createErr)
	}
}

func (brService *BeerServiceImpl) Delete(beerID uint, userID uint) {
	beerToDelete, err := brService.BeerRepository.GetBeerById(beerID)
	switch {
	case err == sql.ErrNoRows:
		panic("The beer is already deleted")
	case err != nil:
		panic(err)
	default:
		if beerToDelete.AuthorId != userID {
			panic("You have to be the author of the beer to be able to delete it")
		}
		delErr := brService.BeerRepository.DeleteBeer(beerID)
		if delErr != nil {
			panic(delErr)
		}
	}
}

func (brService *BeerServiceImpl) FindAll() []model.BeerCompact {
	result, err := brService.BeerRepository.GetAllBeers()
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return result
	}
}

func (brService *BeerServiceImpl) FindById(beerID uint) *model.BeerDetailed {
	result, err := brService.BeerRepository.GetBeerById(beerID)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		var comments []model.Comment
		var cmntErr error
		comments, cmntErr = brService.CommentRepository.GetAllBeerComments(beerID)

		if cmntErr != nil {
			comments = []model.Comment{}
		}
		beerResponse := model.BeerDetailed{
			ID:           result.ID,
			Name:         result.Name,
			Description:  result.Description,
			Thumbnail:    result.Thumbnail,
			CommentCount: result.CommentCount,
			UpvoteCount:  result.UpvoteCount,
			AuthorId:     result.AuthorId,
			Comments:     comments,
		}
		return &beerResponse
	}
}

func (brService *BeerServiceImpl) FindByUser(userID uint) []model.BeerCompact {
	result, err := brService.BeerRepository.GetBeersByUser(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return result
	}
}

func (brService *BeerServiceImpl) Update(beerID uint, beer model.BeerMutate) {
	beerToUpdate, err := brService.BeerRepository.GetBeerById(beerID)
	switch {
	case err == sql.ErrNoRows:
		panic("The beer to update was not found")
	case err != nil:
		panic(err)
	default:
		// TODO: determine what does it do
		valErr := brService.Validate.Struct(beer)
		if valErr != nil {
			panic(valErr)
		}
		if beerToUpdate.AuthorId != beer.AuthorId {
			panic("You have to be the author of the beer to be able to update it")
		}
		brService.BeerRepository.UpdateBeer(beerID, beer)
	}
}
