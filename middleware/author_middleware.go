package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
            c.Abort()
            return
        }

        if userRole == nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Role not defined."})
            c.Abort()
            return
        }
        roleStr, ok := userRole.(string)
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
            c.Abort()
            return
        }

        if roleStr != "author" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Author role required."})
            c.Abort()
            return
        }

        c.Next()
    }
}