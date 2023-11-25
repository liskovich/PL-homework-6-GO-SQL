package model

// TODO: change created date type to date

type Comment struct {
	ID          uint   `json:"id"`
	AuthorID    uint   `json:"user"`
	Content     string `json:"content"`
	CreatedDate string `json:"createdDate"`
	BeerID      uint   `json:"beer"`
}

type CommentMutate struct {
	AuthorID    uint   `json:"user"`
	Content     string `json:"content"`
	CreatedDate string `json:"createdDate"`
	BeerID      uint   `json:"beer"`
}
