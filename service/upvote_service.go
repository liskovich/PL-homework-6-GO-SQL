package service

import (
	"example.com/api/db"
	"example.com/api/model"
	"github.com/go-playground/validator/v10"
)

type UpvoteService interface {
	Upvote(upvote model.Upvote)
	Downvote(downvote model.Upvote)
	CheckIfUserUpvoted(upvote model.Upvote) bool
}

type UpvoteServiceImpl struct {
	UpvoteRepository db.UpvoteRepository
	Validate         *validator.Validate
}

func NewUpvoteService(
	upvoteRepository db.UpvoteRepository,
	validate *validator.Validate) UpvoteService {
	return &UpvoteServiceImpl{
		UpvoteRepository: upvoteRepository,
		Validate:         validate,
	}
}

func (upvtService *UpvoteServiceImpl) Upvote(upvote model.Upvote) {
	valErr := upvtService.Validate.Struct(upvote)
	if valErr != nil {
		panic(valErr)
	}
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(upvote)
	if exists || err != nil {
		panic(err)
	}
	upvoteErr := upvtService.UpvoteRepository.CreateUpvote(upvote)
	if upvoteErr != nil {
		panic(upvoteErr)
	}
}

func (upvtService *UpvoteServiceImpl) Downvote(downvote model.Upvote) {
	valErr := upvtService.Validate.Struct(downvote)
	if valErr != nil {
		panic(valErr)
	}
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(downvote)
	if exists || err != nil {
		panic(err)
	}
	downvoteErr := upvtService.UpvoteRepository.DeleteUpvote(downvote)
	if downvoteErr != nil {
		panic(downvoteErr)
	}
}

func (upvtService *UpvoteServiceImpl) CheckIfUserUpvoted(upvote model.Upvote) bool {
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(upvote)
	if err != nil {
		panic(err)
	}
	return exists
}
