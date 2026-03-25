package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"huynhmanh.com/gin/internal/util"
)

func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized,gin.H{"error":"missing token"})
			c.Abort()
			return
		}
		

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized,gin.H{"message":"invalid token format"})
			c.Abort()
			return 
		}
		
		token := parts[1]

		claims,err := util.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid token"})
			c.Abort()
			return
		}

		c.Set("userID",claims.UserID)
		c.Set("email",claims.Email)

		c.Next()
	}
}