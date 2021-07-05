package dtos

type CommentDTO struct {
	Text     string `json:"text"`
	Date     string `json:"date"`
	Username string `json:"username"`
}
