package db

import (
	"database/sql"

	"example.com/api/model"
)

type UpvoteRepository interface {
	CreateUpvote(upvote model.Upvote) error
	DeleteUpvote(upvote model.Upvote) error
	CheckUpvoteExists(upvote model.Upvote) (bool, error)
}

type upvoteRepo struct {
	db *sql.DB
}

func NewUpvoteRepository(db *sql.DB) UpvoteRepository {
	return &upvoteRepo{db: db}
}

func (upvRepo *upvoteRepo) CreateUpvote(upvote model.Upvote) error {
	_, err := upvRepo.db.Exec(InsertUpvoteQuery, upvote.UserID, upvote.BeerID)
	return err
}

func (upvRepo *upvoteRepo) DeleteUpvote(upvote model.Upvote) error {
	_, err := upvRepo.db.Exec(DeleteUpvoteQuery, upvote.UserID, upvote.BeerID)
	return err
}

func (upvRepo *upvoteRepo) CheckUpvoteExists(upvote model.Upvote) (bool, error) {
	var count int
	err := upvRepo.db.QueryRow(CheckIfUpvoteExistsQuery, upvote.UserID, upvote.BeerID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
