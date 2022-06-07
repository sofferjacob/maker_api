package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/tracking"
)

func PostEvent(c *gin.Context) {
	event := tracking.Event{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	event.Uid = uid
	err := event.Send()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
