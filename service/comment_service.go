package service

import (
	"database/sql"

	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type CommentService interface {
	Create(comment model.CommentMutate)
	FindAllUsersComments(userID uint) []model.Comment
}

type CommentServiceImpl struct {
	CommentRepository db.CommentRepository
	BeerRepository    db.BeerRepository
	Validate          *validator.Validate
}

func NewCommentService(
	commentRepository db.CommentRepository,
	beerRepository db.BeerRepository,
	validate *validator.Validate,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository: commentRepository,
		BeerRepository:    beerRepository,
		Validate:          validate,
	}
}

func (cmntService *CommentServiceImpl) Create(comment model.CommentMutate) {
	_, beerErr := cmntService.BeerRepository.GetBeerById(comment.BeerID)
	switch {
	case beerErr == sql.ErrNoRows:
		panic("The beer to comment was not found")
	case beerErr != nil:
		panic(beerErr)
	default:
		// TODO: determine what does it do
		err := cmntService.Validate.Struct(comment)
		if err != nil {
			panic(err)
		}
		cmntErr := cmntService.CommentRepository.CreateComment(comment)
		if cmntErr != nil {
			panic(cmntErr)
		}
	}
}

func (cmntService *CommentServiceImpl) FindAllUsersComments(userID uint) []model.Comment {
	result, err := cmntService.CommentRepository.GetAllUsersComments(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err)
	default:
		return result
	}
}
