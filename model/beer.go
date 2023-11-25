package model

type Beer struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}
