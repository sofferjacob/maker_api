package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/models"
	"github.com/sofferjacob/maker_api/tracking"
)

func CreateLevel(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	level := models.Level{}
	if err := c.ShouldBindJSON(&level); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	level.Uid = uid
	id, err := level.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": id})
	event := tracking.Event{
		EventType: "level_create",
		Uid:       uid,
		LevelId:   id,
	}
	event.Send()
}

type CreateFromDraftParams struct {
	DraftId     int    `json:"draftId" binding:"required"`
	Name        string `json:"name"`
	Difficulty  int    `json:"difficulty" binding:"required"`
	Description string `json:"description" binding:"required"`
	Theme       int    `json:"theme" binding:"required"`
}

func CreateLevelFromDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := CreateFromDraftParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	draft := models.Draft{Id: params.DraftId}
	err := draft.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if draft.Uid != uid {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	if draft.CourseData == nil || len(draft.CourseData) == 0 {
		c.JSON(400, gin.H{"error": "no course data"})
		return
	}
	level := models.Level{
		Name:        params.Name,
		Description: params.Description,
		Difficulty:  params.Difficulty,
		Theme:       params.Theme,
		Uid:         uid,
	}
	id, err := level.CreateFromDraft(draft)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": id})
	event := tracking.Event{
		EventType: "level_create",
		Uid:       uid,
		LevelId:   id,
		DraftId:   params.DraftId,
	}
	event.Send()
}

func GetLevel(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil || id == 0 || param == "" {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	level := models.Level{Id: id}
	err = level.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "level": level})
}

func GetLevelInfo(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil || id == 0 || param == "" {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	level := models.Level{Id: id}
	err = level.GetInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "level": level})
}

type UpdateLevelParams struct {
	Id          int                    `json:"id" binding:"required"`
	Name        string                 `json:"name"`
	Difficulty  int                    `json:"difficulty"`
	Description string                 `json:"description"`
	Theme       int                    `json:"theme"`
	CourseData  map[string]interface{} `json:"courseData"`
}

func UpdateLevel(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := UpdateLevelParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	level := models.Level{
		Id:          params.Id,
		Uid:         uid,
		Name:        params.Name,
		Difficulty:  params.Difficulty,
		Description: params.Description,
		Theme:       params.Theme,
		CourseData:  params.CourseData,
	}
	err := level.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
	event := tracking.Event{
		EventType: "level_update",
		Uid:       uid,
		LevelId:   params.Id,
	}
	event.Send()
}

type UpdateFromDraftParams struct {
	LevelId int `json:"levelId" binding:"required"`
	DraftId int `json:"draftId" binding:"required"`
}

func UpdateLevelFromDraft(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := UpdateFromDraftParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	draft := models.Draft{Id: params.DraftId}
	err := draft.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if draft.LevelId != params.LevelId || draft.Uid != uid {
		c.JSON(400, gin.H{"error": "invalid draft"})
		return
	}
	err = models.UpdateLevelFomDraft(params.LevelId, draft)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
	event := tracking.Event{
		EventType: "level_update",
		Uid:       uid,
		LevelId:   params.LevelId,
	}
	event.Send()
}

func DeleteLevel(c *gin.Context) {
	c.JSON(500, gin.H{"error": "not implemented", "message": "coming soon"})
}

type QueryFTSParams struct {
	Query string `json:"query" binding:"required"`
}

func QueryLevels(c *gin.Context) {
	params := QueryFTSParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res, err := models.QueryLevelFTS(params.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "results": res})
}

func TrendingLevels(c *gin.Context) {
	levels, err := models.TrendingLevels()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "levels": levels})
}

func Leaderboard(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.Atoi(p)
	if id == 0 || err != nil || p == "" {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := models.GetLeaderboard(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "result": res})
}

func GetUserLevels(c *gin.Context) {
	p := c.Param("uid")
	uid, err := strconv.Atoi(p)
	if err != nil || uid == 0 {
		c.JSON(400, gin.H{"error": "invalid uid"})
		return
	}
	levels, err := models.GetUserLevels(uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "levels": levels})
}
