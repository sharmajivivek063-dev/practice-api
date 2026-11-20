package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"practice-api/connection"
	"practice-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register User API
func RegisterUser(c *gin.Context) {

	var req models.RegisterRequest

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Password hashing failed",
		})
		return
	}

	// 3. Insert into database
	query := `
        INSERT INTO users (name, email, password, mobile)
        VALUES (?, ?, ?, ?)
    `

	result, err := connection.DB.Exec(
		query,
		req.Name,
		req.Email,
		string(hashedPassword),
		req.Mobile,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to register user",
			"error":   err.Error(),
		})
		return
	}

	// 4. Get inserted ID
	lastID, err := result.LastInsertId()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user ID",
		})
		return
	}

	// 5. Prepare response
	user := models.UserResponse{
		ID:     int(lastID),
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
	}

	// 6. Send response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    user,
	})
}

// Login User API
func LoginUser(c *gin.Context) {

	var req models.LoginRequest

	// 1. Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// 2. Find user in database
	var userID int
	var hashedPassword string

	query := `SELECT id, password FROM users WHERE email = ?`

	err := connection.DB.QueryRow(
		query,
		req.Email,
	).Scan(&userID, &hashedPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// 3. Compare password
	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// 4. Get JWT secret
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "JWT secret is not configured",
		})
		return
	}

	// 5. Get expiry time
	expiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))

	if err != nil || expiry <= 0 {
		expiry = 3600
	}

	expirationTime := time.Now().Add(
		time.Duration(expiry) * time.Second,
	)

	// 6. Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   req.Email,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// 7. Create token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Token generation failed",
		})
		return
	}

	// 8. Send response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"access_token": tokenString,
			"token_type":   "Bearer",
			"expires_in":   expiry,
			"user_id":      userID,
		},
	})
}
