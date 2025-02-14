package models

import "time"

// Follow represents the follower-followee relationship.
type Follow struct {
	FollowerID int       `json:"follower_id"`
	FolloweeID int       `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
}
