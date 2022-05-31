package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func getClaims(c *gin.Context) *jwt.StandardClaims {
	claimsObj, ok := c.Get("user-claims")
	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
		return &jwt.StandardClaims{}
	}
	claims, ok := claimsObj.(*jwt.StandardClaims)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "invalid claims"})
		return &jwt.StandardClaims{}
	}
	return claims
}
