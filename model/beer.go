package model

type BeerCompact struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	CommentCount uint   `json:"comment_count"`
	UpvoteCount  uint   `json:"upvote_count"`
	AuthorId     uint   `json:"author_id"`
}

type BeerDetailed struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Thumbnail    string    `json:"thumbnail"`
	CommentCount uint      `json:"comment_count"`
	Comments     []Comment `json:"comments"`
	UpvoteCount  uint      `json:"upvote_count"`
	AuthorId     uint      `json:"author_id"`
}

type BeerMutate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	AuthorId    uint   `json:"author_id"`
}
