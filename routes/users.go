package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/models"
)

func GetUser(c *gin.Context) {
	p := c.Param("id")
	if p == "" {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	u, err := models.GetUser(p)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "user": u})
}

func QueryUsers(c *gin.Context) {
	params := QueryFTSParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	users, err := models.QueryUserFTS(params.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "results": users})
}
