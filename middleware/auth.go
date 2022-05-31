package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const BEARER_SCHEMA string = "Bearer"

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < len(BEARER_SCHEMA)+10 {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden", "message": "a valid token must be provided"})
			return
		}
		token := authHeader[len(BEARER_SCHEMA)+1:]
		tk, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("AUTH_KEY")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden", "message": err.Error()})
			return
		}
		claims, ok := tk.Claims.(*jwt.StandardClaims)
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden", "message": "claim extraction failed"})
			return
		}
		if claims.ExpiresAt < time.Now().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden", "message": "expired token"})
			return
		}
		c.Set("user-claims", claims)
	}
}
