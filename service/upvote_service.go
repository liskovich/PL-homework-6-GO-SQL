package service

import (
	"database/sql"

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
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(upvote)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		panic(err)
	default:
		if !exists {
			valErr := upvtService.Validate.Struct(upvote)
			if valErr != nil {
				panic(valErr)
			}
			upvoteErr := upvtService.UpvoteRepository.CreateUpvote(upvote)
			if upvoteErr != nil {
				panic(upvoteErr)
			}
		} else {
			panic("You already upvoted this beer")
		}
	}
}

func (upvtService *UpvoteServiceImpl) Downvote(downvote model.Upvote) {
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(downvote)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		panic(err)
	default:
		if exists {
			valErr := upvtService.Validate.Struct(downvote)
			if valErr != nil {
				panic(valErr)
			}
			downvoteErr := upvtService.UpvoteRepository.DeleteUpvote(downvote)
			if downvoteErr != nil {
				panic(downvoteErr)
			}
		} else {
			panic("This beer is not upvoted")
		}
	}
}

func (upvtService *UpvoteServiceImpl) CheckIfUserUpvoted(upvote model.Upvote) bool {
	exists, err := upvtService.UpvoteRepository.CheckUpvoteExists(upvote)
	if err != nil {
		panic(err)
	}
	return exists
}
