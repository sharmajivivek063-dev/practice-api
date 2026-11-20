package routes

import (
	"practice-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {

	auth := router.Group("/api/v1/auth")

	auth.POST("/register", controllers.RegisterUser)
	auth.POST("/login", controllers.LoginUser)

}
