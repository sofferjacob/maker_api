package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/sofferjacob/maker_api/db"
)

type Level struct {
	Id          int                    `db:"id" json:"id"`
	Difficulty  int                    `db:"difficulty" json:"difficulty" binding:"required"`
	Name        string                 `db:"name" json:"name" binding:"required"`
	Description string                 `db:"description" json:"description" binding:"required"`
	Uid         int                    `db:"uid" json:"uid"`
	Created     time.Time              `db:"created" json:"created"`
	Updated     sql.NullTime           `db:"updated" json:"updated"`
	Theme       int                    `db:"theme" json:"theme" binding:"required"`
	Car         int                    `db:"car" json:"car" binding:"required"`
	Soundtrack  int                    `db:"soundtrack" json:"soundtrack" binding:"required"`
	CourseData  map[string]interface{} `db:"course_data" json:"courseData" binding:"required"`
}

type DBLevel struct {
	Id          int             `db:"id" json:"id"`
	Difficulty  int             `db:"difficulty" json:"difficulty" binding:"required"`
	Name        string          `db:"name" json:"name" binding:"required"`
	Description string          `db:"description" json:"description" binding:"required"`
	Uid         int             `db:"uid" json:"uid"`
	Created     time.Time       `db:"created" json:"created"`
	Updated     sql.NullTime    `db:"updated" json:"updated"`
	Theme       int             `db:"theme" json:"theme" binding:"required"`
	Ts          string          `db:"ts"`
	Car         int             `db:"car" json:"car" binding:"required"`
	Soundtrack  int             `db:"soundtrack" json:"soundtrack" binding:"required"`
	CourseData  json.RawMessage `db:"course_data" json:"courseData" binding:"required"`
}

func (db *DBLevel) ToLevel(l *Level) {
	cd := map[string]interface{}{}
	json.Unmarshal(db.CourseData, &cd)
	l.Id = db.Id
	l.Difficulty = db.Difficulty
	l.Name = db.Name
	l.Description = db.Description
	l.Uid = db.Uid
	l.Created = db.Created
	l.Updated = db.Updated
	l.Theme = db.Theme
	l.Car = db.Car
	l.Soundtrack = db.Soundtrack
	l.CourseData = cd
}

func (l *Level) Create() (int, error) {
	if l.Difficulty == 0 || l.Name == "" || l.Description == "" || l.Uid == 0 || l.Theme == 0 || l.CourseData == nil || l.Car == 0 || l.Soundtrack == 0 {
		return -1, errors.New("missing required struct fields")
	}
	query := "INSERT INTO levels (difficulty, name, description, uid, theme, car, soundtrack) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := db.Client.Client.Get(&id, query, l.Difficulty, l.Name, l.Description, l.Uid, l.Theme, l.Car, l.Soundtrack)
	if err != nil {
		return -1, err
	}
	courseData := CourseData{LevelId: id, MapData: l.CourseData}
	err = courseData.Create()
	if err != nil {
		db.Client.Client.Exec("DELETE FROM levels WHERE id = $1;", id)
		return -1, err
	}
	return id, nil
}

func (l *Level) CreateFromDraft(d Draft) (int, error) {
	var name string
	if l.Name != "" {
		name = l.Name
	} else if d.Name != "" {
		name = d.Name
	} else {
		return -1, errors.New("missing required field name")
	}
	if l.Difficulty == 0 || l.Description == "" || l.Uid == 0 || l.Theme == 0 || d.CourseData == nil {
		return -1, errors.New("missing required struct fields")
	}
	query := "INSERT INTO levels (difficulty, name, description, uid, theme, car, soundtrack) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := db.Client.Client.Get(&id, query, l.Difficulty, name, l.Description, l.Uid, l.Theme, d.Car, d.Soundtrack)
	if err != nil {
		return -1, err
	}
	courseData := CourseData{LevelId: id, MapData: d.CourseData}
	err = courseData.Create()
	if err != nil {
		db.Client.Client.Exec("DELETE FROM levels WHERE id = $1;", id)
		return -1, err
	}
	d.Delete()
	return id, nil
}

func (l *Level) Get() error {
	if l.Id == 0 {
		return errors.New("missing required field Id")
	}
	res := DBLevel{}
	query := "SELECT l.*, c.map_data course_data FROM levels l INNER JOIN course_data c ON l.id = c.level_id WHERE l.id = $1;"
	err := db.Client.Client.Get(&res, query, l.Id)
	res.ToLevel(l)
	return err
}

func (l *Level) GetInfo() error {
	if l.Id == 0 {
		return errors.New("missing required field Id")
	}
	d := DBLevel{}
	query := "SELECT * FROM levels WHERE id = $1;"
	err := db.Client.Client.Get(&d, query, l.Id)
	d.ToLevel(l)
	return err
}

func GetUserLevels(uid int) ([]Level, error) {
	query := "SELECT * FROM levels WHERE uid = $1;"
	res := []DBLevel{}
	err := db.Client.Client.Select(&res, query, uid)
	levels := make([]Level, len(res))
	for i, v := range res {
		v.ToLevel(&levels[i])
	}
	return levels, err
}

func (l *Level) Update() error {
	if l.Id == 0 || l.Uid == 0 {
		return errors.New("missing required fields Id, Uid")
	}
	query := db.Update("levels")
	if l.Name != "" {
		query = query.Set("name", l.Name)
	}
	if l.Difficulty != 0 {
		query = query.Set("difficulty", l.Difficulty)
	}
	if l.Description != "" {
		query = query.Set("description", l.Description)
	}
	if l.Theme != 0 {
		query = query.Set("theme", l.Theme)
	}
	queryStr, args := query.Where("id", "=", l.Id).
		And("uid", "=", l.Uid).Query()
	_, err := db.Client.Client.Exec(queryStr, args...)
	if err != nil {
		return err
	}
	if l.CourseData != nil {
		cd := CourseData{LevelId: l.Id, MapData: l.CourseData}
		err = cd.Update()
	}
	return err
}

func UpdateLevelFomDraft(levelId int, draft Draft) error {
	levelQuery := "SELECT uid FROM levels WHERE id = $1;"
	var levelUid int
	err := db.Client.Client.Get(&levelUid, levelQuery, levelId)
	if err != nil {
		return err
	}
	if levelUid != draft.Uid || draft.CourseData == nil {
		return errors.New("invalid draft")
	}
	level := Level{Id: levelId, Uid: levelUid, Name: draft.Name, CourseData: draft.CourseData, Car: draft.Car, Soundtrack: draft.Soundtrack}
	return level.Update()
}

func QueryLevelFTS(query string) ([]Level, error) {
	res := []DBLevel{}
	err := db.Client.Client.Select(&res, "SELECT * FROM query_gin(null::levels, $1);", query)
	if err != nil {
		return nil, err
	}
	levels := make([]Level, 0, len(res))
	for _, v := range res {
		l := Level{}
		v.ToLevel(&l)
		levels = append(levels, l)
	}
	return levels, err
}

func DeleteLevel(levelId int) error {
	// TODO, maybe with a procedure
	return errors.New("not implemented")
}

func TrendingLevels() ([]Level, error) {
	query := "SELECT id, difficulty, name, description, uid, created, updated, theme FROM trending_levels;"
	res := []Level{}
	err := db.Client.Client.Select(&res, query)
	return res, err
}
