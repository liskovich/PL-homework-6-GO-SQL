package model

type BeerCompact struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Thumbnail     string `json:"thumbnail"`
	CommentsCount uint   `json:"comments_count"`
	UpvotesCount  uint   `json:"upvotes_count"`
}

type BeerDetailed struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Thumbnail     string    `json:"thumbnail"`
	CommentsCount uint      `json:"comments_count"`
	Comments      []Comment `json:"comments"`
	UpvotesCount  uint      `json:"upvotes_count"`
}

type BeerMutate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
