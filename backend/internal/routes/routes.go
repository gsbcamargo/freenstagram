package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsbcamargo/freenstagram/backend/internal/handlers"
)

// SetupRouter configures the API routes.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public authentication routes.
	router.POST("/signup", handlers.SignUp)
	router.POST("/login", handlers.Login)

	// Protected API routes.
	api := router.Group("/api")
	{
		api.POST("/posts", handlers.CreatePost)
		api.GET("/posts", handlers.GetPosts)
		api.POST("/follow/:username", handlers.FollowUser)
		api.DELETE("/unfollow/:username", handlers.UnfollowUser)
	}

	return router
}
