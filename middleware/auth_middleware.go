package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization token is required",
			})
			c.Abort()
			return
		}

		// 2. Check Bearer token format
		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Get JWT secret
		secret := os.Getenv("JWT_SECRET")

		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "JWT secret is not configured",
			})
			c.Abort()
			return
		}

		// 4. Parse and verify token
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// Accept only HS256
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 5. Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)

		if !ok || userID <= 0 || userID != float64(int64(userID)) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		// 6. Store verified user ID in Gin context
		c.Set("user_id", int64(userID))

		// 7. Continue to controller
		c.Next()
	}
}
