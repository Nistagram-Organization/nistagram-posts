package dtos

type PostDTO struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Image       string `json:"image"`
	Username    string `json:"username"`
	Liked       bool   `json:"liked"`
	Disliked    bool   `json:"disliked"`
	Favorited   bool   `json:"favorited"`
	Likes       uint   `json:"likes"`
	Dislikes    uint   `json:"dislikes"`
	Comments    []CommentDTO
}
