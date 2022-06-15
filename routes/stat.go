package routes

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/db"
)

type LevelStartsResult struct {
	GameStarts int       `db:"game_starts" json:"gameStarts"`
	Date       time.Time `db:"date" json:"date"`
}

type GetLevelStartsParams struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
	Gt   int       `json:"gt"`
	Lt   int       `json:"lt"`
}

type AvgTimeParams struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
	Gt   float64   `json:"gt"`
	Lt   float64   `json:"lt"`
}

func GetLevelStarts(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	params := GetLevelStartsParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	qb := db.SelectFrom("events").Select("COUNT(*) game_starts").Select("date(timestamp)").Where("event_type", "=", "game_start").
		And("level_id", "=", id)
	base := time.Time{}
	if params.From != base {
		qb = qb.And("timestamp", ">=", params.From)
	}
	if params.To != base {
		qb = qb.And("timestamp", "<=", params.To)
	}
	query, args := qb.GroupBy("date(timestamp)").Query()
	//query, args := "SELECT COUNT(*) game_starts, date(timestamp) FROM events WHERE event_type = 'game_start' AND level_id = $1 GROUP BY date(timestamp);"
	res := []LevelStartsResult{}
	err = db.Client.Client.Select(&res, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if params.Gt != 0 || params.Lt != 0 {
		for i := 0; i < len(res); i++ {
			if (params.Gt != 0 && res[i].GameStarts < params.Gt) || (params.Lt != 0 && res[i].GameStarts > params.Lt) {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			}
		}
	}
	c.JSON(200, gin.H{"status": "ok", "result": res})
}

type LevelCompleteResult struct {
	GameComplete int       `db:"game_complete" json:"gameCompletes"`
	Date         time.Time `db:"date" json:"date"`
}

func GetLevelCompletes(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	params := GetLevelStartsParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	qb := db.SelectFrom("events").Select("COUNT(*) game_complete").Select("date(timestamp)").Where("event_type", "=", "game_finish").
		And("level_id", "=", id)
	base := time.Time{}
	if params.From != base {
		qb = qb.And("timestamp", ">=", params.From)
	}
	if params.To != base {
		qb = qb.And("timestamp", "<=", params.To)
	}
	query, args := qb.GroupBy("date(timestamp)").Query()
	//query, args := "SELECT COUNT(*) game_starts, date(timestamp) FROM events WHERE event_type = 'game_start' AND level_id = $1 GROUP BY date(timestamp);"
	res := []LevelCompleteResult{}
	err = db.Client.Client.Select(&res, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if params.Gt != 0 || params.Lt != 0 {
		for i := 0; i < len(res); i++ {
			if (params.Gt != 0 && res[i].GameComplete < params.Gt) || (params.Lt != 0 && res[i].GameComplete > params.Lt) {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			}
		}
	}
	c.JSON(200, gin.H{"status": "ok", "result": res})
}

type AvgTimeResult struct {
	AvgTime float64   `db:"avg_time" json:"avgTime"`
	Date    time.Time `db:"date" json:"date"`
}

func GetAvgTime(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	params := AvgTimeParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	qb := db.SelectFrom("events").Select("AVG(time) avg_time").Select("date(timestamp)").Where("event_type", "=", "game_finish").
		And("level_id", "=", id)
	base := time.Time{}
	if params.From != base {
		qb = qb.And("timestamp", ">=", params.From)
	}
	if params.To != base {
		qb = qb.And("timestamp", "<=", params.To)
	}
	query, args := qb.GroupBy("date(timestamp)").Query()
	//query, args := "SELECT COUNT(*) game_starts, date(timestamp) FROM events WHERE event_type = 'game_start' AND level_id = $1 GROUP BY date(timestamp);"
	res := []AvgTimeResult{}
	err = db.Client.Client.Select(&res, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if params.Gt != 0 || params.Lt != 0 {
		for i := 0; i < len(res); i++ {
			if (params.Gt != 0 && res[i].AvgTime < params.Gt) || (params.Lt != 0 && res[i].AvgTime > params.Lt) {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			}
		}
	}
	c.JSON(200, gin.H{"status": "ok", "result": res})
}

type UniqueUsersResult struct {
	UniqueUsers int `db:"unique_users" json:"uniqueUsers"`
	Month       int `db:"month" json:"month"`
}

func GetUniqueUsers(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	params := GetLevelStartsParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err != nil || id == 0 {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	qb := db.SelectFrom("events").Select("COUNT(DISTINCT uid) unique_users").Select("EXTRACT(MONTH FROM timestamp) AS month").Where("event_type", "=", "game_finish").
		And("level_id", "=", id)
	base := time.Time{}
	if params.From != base {
		qb = qb.And("timestamp", ">=", params.From)
	}
	if params.To != base {
		qb = qb.And("timestamp", "<=", params.To)
	}
	query, args := qb.GroupBy("EXTRACT(MONTH FROM timestamp)").Query()
	//query, args := "SELECT COUNT(*) game_starts, date(timestamp) FROM events WHERE event_type = 'game_start' AND level_id = $1 GROUP BY date(timestamp);"
	res := []UniqueUsersResult{}
	err = db.Client.Client.Select(&res, query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if params.Gt != 0 || params.Lt != 0 {
		for i := 0; i < len(res); i++ {
			if (params.Gt != 0 && res[i].UniqueUsers < params.Gt) || (params.Lt != 0 && res[i].UniqueUsers > params.Lt) {
				res[i] = res[len(res)-1]
				res = res[:len(res)-1]
			}
		}
	}
	c.JSON(200, gin.H{"status": "ok", "result": res})
}
