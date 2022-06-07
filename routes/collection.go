package routes

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/models"
)

type CreateCollectionParams struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateCollection(c *gin.Context) {
	claims := getClaims(c)
	params := CreateCollectionParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	uid, _ := strconv.Atoi(claims.Subject)
	collection := models.Collection{
		Name:        params.Name,
		Description: params.Description,
		Uid:         uid,
	}
	id, err := collection.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "id": id})
}

func GetCollection(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if param == "" || err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := models.GetCollection(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "collection": res})
}

func GetUserCollections(c *gin.Context) {
	param := c.Param("uid")
	uid, err := strconv.Atoi(param)
	if param == "" || err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	res, err := models.GetUserCollections(uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "collections": res})
}

func DeleteCollection(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if param == "" || err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	collection := models.Collection{
		Id:  id,
		Uid: uid,
	}
	err = collection.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

type UpdateCollectionParams struct {
	Id          int    `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func UpdateCollection(c *gin.Context) {
	claims := getClaims(c)
	params := UpdateCollectionParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	uid, _ := strconv.Atoi(claims.Subject)
	collection := models.Collection{
		Id:          params.Id,
		Uid:         uid,
		Name:        params.Name,
		Description: params.Description,
	}
	err := collection.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func QueryCollections(c *gin.Context) {
	params := QueryFTSParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	col, err := models.QueryCollectionsFTS(params.Query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "results": col})
}

type LinkLevelParams struct {
	LevelId      int `json:"levelId" binding:"required"`
	CollectionId int `json:"collectionId" binding:"required"`
}

func LinkLevel(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := LinkLevelParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	isOwn, err := models.IsOwnCollection(params.CollectionId, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !isOwn {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	obj := models.CollectionLevels{
		LevelId:      params.LevelId,
		CollectionId: params.CollectionId,
	}
	err = obj.Create()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func UnlinkLevel(c *gin.Context) {
	claims := getClaims(c)
	uid, _ := strconv.Atoi(claims.Subject)
	params := LinkLevelParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	isOwn, err := models.IsOwnCollection(params.CollectionId, uid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !isOwn {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	obj := models.CollectionLevels{
		LevelId:      params.LevelId,
		CollectionId: params.CollectionId,
	}
	err = obj.Delete()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}

func GetCollectionLevels(c *gin.Context) {
	p := c.Param("id")
	collecionId, err := strconv.Atoi(p)
	if err != nil || p == "" {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	levels, err := models.GetCollectionLevels(collecionId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "levels": levels})
}

func TrendingCollections(c *gin.Context) {
	cls, err := models.TrendingCollections()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "collections": cls})
}
