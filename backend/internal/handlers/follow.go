package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsbcamargo/freenstagram/backend/internal/db"
)

// FollowUser allows a user to follow another user by username.
func FollowUser(c *gin.Context) {
	// For simplicity, assume the user is authenticated and use a hard-coded followerID.
	// Replace this with the authenticated user ID from your middleware.
	followerID := 1

	followeeUsername := c.Param("username")
	var followeeID int
	err := db.DB.QueryRow(
		`SELECT id FROM users WHERE username = $1`,
		followeeUsername,
	).Scan(&followeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}

	_, err = db.DB.Exec(
		`INSERT INTO follows (follower_id, followee_id, created_at)
         VALUES ($1, $2, $3)`,
		followerID, followeeID, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user followed successfully"})
}

// UnfollowUser allows a user to unfollow another user by username.
func UnfollowUser(c *gin.Context) {
	// For simplicity, assume the user is authenticated and use a hard-coded followerID.
	followerID := 1

	followeeUsername := c.Param("username")
	var followeeID int
	err := db.DB.QueryRow(
		`SELECT id FROM users WHERE username = $1`,
		followeeUsername,
	).Scan(&followeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}

	_, err = db.DB.Exec(
		`DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`,
		followerID, followeeID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user unfollowed successfully"})
}
