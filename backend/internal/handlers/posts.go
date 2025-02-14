package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsbcamargo/freenstagram/backend/internal/db"
)

const maxImages = 4

// CreatePost handles uploading a new post with up to 4 images.
func CreatePost(c *gin.Context) {
	// For simplicity, assume the user is authenticated and userID is available.
	// In production, extract the userID from the authentication middleware.
	userID := 1 // Replace with the real authenticated user ID

	caption := c.PostForm("caption")
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 || len(files) > maxImages {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please upload between 1 and 4 images"})
		return
	}

	// Begin a transaction.
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer tx.Rollback()

	// Insert the post.
	var postID int
	err = tx.QueryRow(
		`INSERT INTO posts (user_id, caption, created_at)
         VALUES ($1, $2, $3) RETURNING id`,
		userID, caption, time.Now(),
	).Scan(&postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	// Process each uploaded image.
	for idx, fileHeader := range files {
		filename := fmt.Sprintf("uploads/%d_%d_%s", postID, idx, fileHeader.Filename)
		if err := c.SaveUploadedFile(fileHeader, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}

		// Optionally, receive width and height parameters.
		widthStr := c.PostForm(fmt.Sprintf("width_%d", idx))
		heightStr := c.PostForm(fmt.Sprintf("height_%d", idx))
		width, _ := strconv.Atoi(widthStr)
		height, _ := strconv.Atoi(heightStr)
		if width == 0 || height == 0 {
			width, height = 800, 600
		}

		_, err = tx.Exec(
			`INSERT INTO post_images (post_id, image_url, width, height)
             VALUES ($1, $2, $3, $4)`,
			postID, filename, width, height,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image info"})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post created successfully", "post_id": postID})
}

// GetPosts retrieves posts in reverse chronological order.
func GetPosts(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT p.id, p.caption, p.created_at, u.username 
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
		LIMIT 50
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch posts"})
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id int
		var caption, username string
		var createdAt time.Time
		if err := rows.Scan(&id, &caption, &createdAt, &username); err != nil {
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id":         id,
			"caption":    caption,
			"username":   username,
			"created_at": createdAt,
		})
	}
	c.JSON(http.StatusOK, posts)
}
