package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsbcamargo/freenstagram/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

// SignUpInput represents the expected JSON payload for user registration.
type SignUpInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignUp handles new user registration.
func SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Insert the new user into the database.
	var userID int
	err = db.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		input.Username, input.Email, string(hashedPassword), time.Now(),
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created", "user_id": userID})
}

// LoginInput represents the expected JSON payload for login.
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user authentication.
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user record.
	var id int
	var passwordHash string
	err := db.DB.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = $1`,
		input.Username,
	).Scan(&id, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}

	// Compare the hashed password.
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// For demonstration, return a dummy token.
	token := "dummy-token-for-user-" + input.Username

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}
