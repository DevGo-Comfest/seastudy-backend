package middleware

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"sea-study/service"
	"sea-study/util"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func UserMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(util.GetJWTSecret()), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["sub"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}
			// Validate UUID format
			userUUID, err := uuid.Parse(userID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid UUID format"})
				c.Abort()
				return
			}

			// Check if the user exists
			err = service.ValidateUserExists(db, userUUID)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}

			c.Set("userID", userID)
			c.Set("userRole", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
