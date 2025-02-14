package models

import "time"

// Post represents a user post.
type Post struct {
	ID        int         `json:"id"`
	UserID    int         `json:"user_id"`
	Caption   string      `json:"caption"`
	CreatedAt time.Time   `json:"created_at"`
	Images    []PostImage `json:"images"`
}

// PostImage represents an image associated with a post.
type PostImage struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	ImageURL string `json:"image_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}
