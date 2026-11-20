package models

import (
	"time"
)

// Registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Mobile   string `json:"mobile" binding:"required,len=10,numeric"`
}

// User response
type UserResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

// Login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// User profile response
type UserProfileResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
}

// update user profile
type UpdateProfileRequest struct {
	Name   string `json:"name" binding:"required"`
	Mobile string `json:"mobile" binding:"required,len=10,numeric"`
}
