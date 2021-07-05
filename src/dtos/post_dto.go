package dtos

type PostDTO struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Timestamp   int64  `json:"timestamp"`
	Image       string `json:"image"`
	Username    string `json:"username"`
	Liked       bool   `json:"liked"`
	Disliked    bool   `json:"disliked"`
	InFavorites bool   `json:"in_favorites"`
	Likes       uint   `json:"likes"`
	Dislikes    uint   `json:"dislikes"`
	Comments    []CommentDTO
}
