package controllers

import (
	"database/sql"
	"net/http"

	"practice-api/connection"
	"practice-api/models"

	"github.com/gin-gonic/gin"
)

// Get User Profile API
func GetProfile(c *gin.Context) {

	// 1. Get verified user ID from middleware
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		return
	}

	userID := userIDValue.(int64)

	// 2. Fetch user profile from database
	var user models.UserProfileResponse

	query := `
        SELECT id, name, email, mobile, created_at
        FROM users
        WHERE id = ?
    `

	err := connection.DB.QueryRow(
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Mobile,
		&user.CreatedAt,
	)

	// 3. Check database errors
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch user profile",
		})
		return
	}

	// 4. Return profile
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// user profile updation controller
func UpdateProfile(c *gin.Context) {

	// 1. Get user ID from JWT middleware
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized access",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// 2. Read request JSON
	var req models.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 3. Update user in database
	query := `
        UPDATE users
        SET name = ?, mobile = ?
        WHERE id = ?
    `

	_, err := connection.DB.Exec(query, req.Name, req.Mobile, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
		})
		return
	}

	// 4. Fetch updated user details
	var user models.UserResponse

	selectQuery := `
        SELECT id, name, email, mobile
        FROM users
        WHERE id = ?
    `

	err = connection.DB.QueryRow(selectQuery, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Mobile,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch updated profile",
		})
		return
	}

	// 5. Send success response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    user,
	})
}
