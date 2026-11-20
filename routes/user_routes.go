package routes

import (
	"practice-api/controllers"
	"practice-api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/api/v1/users")

	user.Use(middleware.AuthMiddleware())

	user.GET("/profile", controllers.GetProfile)

	// Update profile
	user.PUT("/profile", controllers.UpdateProfile)
}
