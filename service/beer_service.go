package service

import (
	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type BeerService interface {
	Create(beer model.BeerMutate)
	Update(beerID uint, beer model.BeerMutate)
	Delete(beerID uint, userID uint)
	FindById(beerID uint) *model.BeerDetailed
	FindByUser(userID uint) []*model.BeerCompact
	FindAll() []*model.BeerCompact
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
	err := brService.Validate.Struct(beer)
	if err != nil {
		panic(err)
	}
	brService.BeerRepository.CreateBeer(beer)
}

func (brService *BeerServiceImpl) Delete(beerID uint, userID uint) {
	beerToDelete, err := brService.BeerRepository.GetBeerById(beerID)
	if err != nil || beerToDelete == nil {
		panic(err)
	}
	if beerToDelete.AuthorId != userID {
		panic("You have to be the author of the beer to be able to delete it")
	}
	delErr := brService.BeerRepository.DeleteBeer(beerID)
	if delErr != nil {
		panic(delErr)
	}
}

func (brService *BeerServiceImpl) FindAll() []*model.BeerCompact {
	result, err := brService.BeerRepository.GetAllBeers()
	if err != nil {
		panic(err)
	}
	return result
}

func (brService *BeerServiceImpl) FindById(beerID uint) *model.BeerDetailed {
	result, err := brService.BeerRepository.GetBeerById(beerID)
	if err != nil {
		panic(err)
	}
	if result != nil {
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
			AuthorId:     result.AuthorId,
			Comments:     comments,
		}
		return &beerResponse
	} else {
		return nil
	}
}

func (brService *BeerServiceImpl) FindByUser(userID uint) []*model.BeerCompact {
	result, err := brService.BeerRepository.GetBeersByUser(userID)
	if err != nil {
		panic(err)
	}
	return result
}

func (brService *BeerServiceImpl) Update(beerID uint, beer model.BeerMutate) {
	beerToUpdate, err := brService.BeerRepository.GetBeerById(beerID)
	if err != nil || beerToUpdate == nil {
		panic(err)
	}
	valErr := brService.Validate.Struct(beer)
	if valErr != nil {
		panic(valErr)
	}
	if beerToUpdate.AuthorId != beer.AuthorId {
		panic("You have to be the author of the beer to be able to update it")
	}
	brService.BeerRepository.UpdateBeer(beerID, beer)
}
