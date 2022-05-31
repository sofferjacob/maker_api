package routes

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/models"
)

type CreateDraftParams struct {
	Name       string                 `json:"name" binding:"required"`
	LevelId    int                    `json:"levelId"`
	CourseData map[string]interface{} `json:"courseData"`
	Theme      int                    `json:"theme"`
}

func CreateDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := CreateDraftParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	draft := models.Draft{
		Name:       params.Name,
		LevelId:    params.LevelId,
		CourseData: params.CourseData,
		Theme:      params.Theme,
		Uid:        uid,
	}
	id, err := draft.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": id})
}

type UpdateDraftParams struct {
	Name       string                 `json:"name"`
	CourseData map[string]interface{} `json:"courseData"`
	Theme      int                    `json:"theme"`
	Id         int                    `json:"id" binding:"required"`
}

func UpdateDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := UpdateDraftParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	draft := models.Draft{
		Id:         params.Id,
		Name:       params.Name,
		CourseData: params.CourseData,
		Theme:      params.Theme,
		Uid:        uid,
	}
	err := draft.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func GetDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	param := c.Param("id")
	if param == "" {
		c.JSON(400, gin.H{"error": "missing required parameter id"})
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	draft := models.Draft{Id: id}
	err = draft.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if draft.Uid != uid {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "draft": draft})
}

func GetDraftFromLevelId(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	param := c.Param("id")
	if param == "" {
		c.JSON(400, gin.H{"error": "missing required parameter id"})
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	draft := models.Draft{LevelId: id}
	err = draft.FromLevelId()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if draft.Uid != uid {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "draft": draft})
}

func DeleteDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	param := c.Param("id")
	if param == "" {
		c.JSON(400, gin.H{"error": "missing required parameter id"})
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	draft := models.Draft{
		Id:  id,
		Uid: uid,
	}
	err = draft.SafeDelete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func GetUserDrafts(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	drafts, err := models.GetUserDrafts(uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "drafts": drafts})
}

func GetLevelDraft(c *gin.Context) {
	fmt.Println("called here")
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	param := c.Param("id")
	levelId, err := strconv.Atoi(param)
	if param == "" || err != nil {
		c.JSON(400, gin.H{"error": "invalid level id"})
		return
	}
	draft, err := models.GetLevelDraft(levelId, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "draft": draft})
}
