package service

import (
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
	Validate          *validator.Validate
}

func NewCommentService(
	commentRepository db.CommentRepository,
	validate *validator.Validate) CommentService {
	return &CommentServiceImpl{
		CommentRepository: commentRepository,
		Validate:          validate,
	}
}

func (cmntService *CommentServiceImpl) Create(comment model.CommentMutate) {
	err := cmntService.Validate.Struct(comment)
	if err != nil {
		panic(err)
	}
	// TODO: check if user is authenticated and beer exists
	cmntErr := cmntService.CommentRepository.CreateComment(comment)
	if cmntErr != nil {
		panic(cmntErr)
	}
}

func (cmntService *CommentServiceImpl) FindAllUsersComments(userID uint) []model.Comment {
	result, err := cmntService.CommentRepository.GetAllUsersComments(userID)
	if err != nil {
		panic(err)
	}
	return result
}
