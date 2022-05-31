package models

import (
	"errors"

	"encoding/json"

	"github.com/sofferjacob/maker_api/db"
)

type CourseData struct {
	Id      int                    `db:"id" json:"id"`
	LevelId int                    `db:"level_id" json:"levelId" binding:"required"`
	MapData map[string]interface{} `db:"map_data" json:"mapData" binding:"required"`
}

func (c *CourseData) Create() error {
	if c.LevelId == 0 {
		return errors.New("missing required field LevelId")
	}
	query := "INSERT INTO course_data (level_id, map_data) VALUES ($1, $2);"
	cd, err := json.Marshal(c.MapData)
	if err != nil {
		return err
	}
	_, err = db.Client.Client.Exec(query, c.LevelId, cd)
	return err
}

func (c *CourseData) Update() error {
	if c.LevelId == 0 {
		return errors.New("missing required field LevelId")
	}
	cd, err := json.Marshal(c.MapData)
	if err != nil {
		return err
	}
	query := "UPDATE course_data SET map_data = $1 WHERE level_id = $2"
	_, err = db.Client.Client.Exec(query, cd, c.LevelId)
	return err
}

func (c *CourseData) Delete() error {
	if c.LevelId == 0 {
		return errors.New("missing required field LevelId")
	}
	query := "DELETE FROM course_data WHERE level_id = $1"
	_, err := db.Client.Client.Exec(query, c.MapData)
	return err
}
