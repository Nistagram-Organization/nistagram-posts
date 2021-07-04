package dtos

type InappropriateContentReportDTO struct {
	AuthorEmail string `json:"author_email"`
	Description string `json:"description"`
	Image       string `json:"image"`
	PostID      uint   `json:"post_id"`
}
