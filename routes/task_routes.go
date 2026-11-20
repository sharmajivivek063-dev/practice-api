package routes

import (
	"practice-api/controllers"
	"practice-api/middleware"

	"github.com/gin-gonic/gin"
)

func TaskRoutes(router *gin.Engine) {

	task := router.Group("/api/v1/tasks")

	task.Use(middleware.AuthMiddleware())

	// Create task
	task.POST("", controllers.CreateTask)

	// List tasks
	task.GET("", controllers.ListTasks)

	// Get single task
	task.GET("/:id", controllers.GetTask)

	// update the task
	task.PUT("/:id", controllers.UpdateTask)

	// Delete the Task
	task.DELETE("/:id", controllers.DeleteTask)

	// get the task Stats
	task.GET("/stats", controllers.GetTaskStats)
}
