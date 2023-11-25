package model

type Upvote struct {
	ID     uint `json:"id"`
	UserID uint `json:"user"`
	BeerID uint `json:"beer"`
}
