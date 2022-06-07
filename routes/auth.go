package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sofferjacob/maker_api/models"
	"github.com/sofferjacob/maker_api/tracking"
)

type LoginParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	params := LoginParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	u := models.User{
		Email:    params.Email,
		Password: params.Password,
	}
	token, err := u.Login()
	if err != nil {
		c.JSON(500, gin.H{"error": "could not authenticate user", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "token": token, "user": u.ToUserData()})
	event := tracking.Event{
		EventType: "user_login",
		Uid:       u.Id,
	}
	event.Send()
}

type RegisterParams struct {
	LoginParams
	Name string `json:"name" binding:"required"`
}

func Register(c *gin.Context) {
	params := RegisterParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	u := models.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}
	err := u.Register()
	if err != nil {
		c.JSON(500, gin.H{"error": "could not register user", "message": err.Error()})
		return
	}
	event := tracking.Event{
		EventType: "user_register",
	}
	event.Send()
	c.JSON(200, gin.H{"status": "ok"})
}

func Profile(c *gin.Context) {
	claimsObj, ok := c.Get("user-claims")
	if !ok {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	claims, ok := claimsObj.(*jwt.StandardClaims)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid claims"})
		return
	}
	u, err := models.GetUser(claims.Subject)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "user": u})
}

type UpdateProfileParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UpdateProfile(c *gin.Context) {
	claimsObj, ok := c.Get("user-claims")
	if !ok {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	claims, ok := claimsObj.(*jwt.StandardClaims)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid claims"})
		return
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}
	params := UpdateProfileParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	u := models.User{
		Id:    id,
		Name:  params.Name,
		Email: params.Email,
	}
	err = u.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
