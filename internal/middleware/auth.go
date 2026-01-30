package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pick-cee/events-api/internal/config"
)

type Claims struct {
	UserID uint   `json:"user_id"`
  Email  string `json:"email"`
  jwt.RegisteredClaims
}

func AuthMidleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()

	}
}

func GetUserId(c *gin.Context) uint {
	userId, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	return userId.(uint)
}

func GetUserEmail(c *gin.Context) string {
	email, exists := c.Get("email")
	if !exists {
		return ""
	}
	return email.(string)
}