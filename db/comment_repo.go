package db

import (
	"database/sql"

	"example.com/api/model"
)

type CommentRepository interface {
	CreateComment(comment model.Comment) error
	GetAllUsersComments(userID uint) ([]*model.Comment, error)
}

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (cmntRepo *commentRepo) CreateComment(comment model.Comment) error {
	_, err := cmntRepo.db.Exec(InsertCommentQuery, comment.AuthorID, comment.Content, comment.CreatedDate, comment.BeerID)
	return err
}

func (cmntRepo *commentRepo) GetAllUsersComments(userID uint) ([]*model.Comment, error) {
	rows, err := cmntRepo.db.Query(SelectAllCommentsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comment{}
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content, &comment.CreatedDate, &comment.BeerID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}
