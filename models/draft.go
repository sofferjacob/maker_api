package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/sofferjacob/maker_api/db"
	"github.com/sofferjacob/maker_api/tracking"
)

type Draft struct {
	Id         int                    `db:"id" json:"id"`
	Name       string                 `db:"name" json:"name"`
	LevelId    int                    `db:"level_id" json:"levelId"`
	Created    time.Time              `db:"created" json:"created"`
	Updated    sql.NullTime           `db:"updated" json:"updated"`
	CourseData map[string]interface{} `db:"course_data" json:"courseData"`
	Theme      int                    `db:"theme" json:"theme"`
	Car        int                    `db:"car" json:"car" binding:"required"`
	Soundtrack int                    `db:"soundtrack" json:"soundtrack" binding:"required"`
	Uid        int                    `db:"uid" json:"uid"`
}

type DbDraft struct {
	Id         int            `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	LevelId    sql.NullInt32  `db:"level_id" json:"levelId"`
	Created    time.Time      `db:"created" json:"created"`
	Updated    sql.NullTime   `db:"updated" json:"updated"`
	CourseData types.JSONText `db:"course_data" json:"courseData"`
	Theme      int            `db:"theme" json:"theme"`
	Car        int            `db:"car" json:"car" binding:"required"`
	Soundtrack int            `db:"soundtrack" json:"soundtrack" binding:"required"`
	Uid        int            `db:"uid" json:"uid"`
}

func (db *DbDraft) LoadToDraft(d *Draft) {
	cd := map[string]interface{}{}
	json.Unmarshal(db.CourseData, &cd)
	d.Id = db.Id
	d.Name = db.Name
	d.LevelId = int(db.LevelId.Int32)
	d.Created = db.Created
	d.Updated = db.Updated
	d.CourseData = cd
	d.Theme = db.Theme
	d.Uid = db.Uid
	d.Car = db.Car
	d.Soundtrack = db.Soundtrack
}

func (d *Draft) CourseDataDb() (types.JSONText, error) {
	return json.Marshal(d.CourseData)
}

func (d *Draft) Create() (int, error) {
	if d.Name == "" || d.Uid == 0 || d.Car == 0 || d.Soundtrack == 0 {
		return -1, errors.New("missing required field name, uid, car, soundtrack")
	}
	qb := db.Insert("drafts").Set("name", d.Name).Set("uid", d.Uid).
		Set("car", d.Car).Set("soundtrack", d.Soundtrack)
	if d.LevelId != 0 {
		qb = qb.Set("level_id", d.LevelId)
	}
	if d.CourseData != nil {
		cd, err := d.CourseDataDb()
		if err != nil {
			return -1, err
		}
		qb = qb.Set("course_data", cd)
	}
	if d.Theme != 0 {
		qb = qb.Set("theme", d.Theme)
	}

	query, args := qb.Returning("id").Query()
	var id int
	err := db.Client.Client.Get(&id, query, args...)
	return id, err
}

func (d *Draft) Update() error {
	if d.Id == 0 || d.Uid == 0 {
		return errors.New("missing required field id")
	}
	qb := db.Update("drafts")
	if d.Name != "" {
		qb = qb.Set("name", d.Name)
	}
	if d.Theme != 0 {
		qb = qb.Set("theme", d.Theme)
	}
	if d.Car != 0 {
		qb = qb.Set("car", d.Car)
	}
	if d.Soundtrack != 0 {
		qb = qb.Set("soundtrack", d.Soundtrack)
	}
	if d.CourseData != nil {
		cdE, err := json.Marshal(d.CourseData)
		cd := types.JSONText(cdE)
		if err != nil {
			return fmt.Errorf("could not encode course data: %v", err.Error())
		}
		fmt.Println(cd)
		qb = qb.Set("course_data", cd)
	}
	query, args := qb.Where("id", "=", d.Id).And("uid", "=", d.Uid).
		Query()
	fmt.Println(query)
	fmt.Println(args)
	_, err := db.Client.Client.Exec(query, args...)
	return err
}

func (d *Draft) Get() error {
	if d.Id == 0 {
		return errors.New("missing required field LevelId")
	}
	res := DbDraft{}
	query := "SELECT * FROM drafts WHERE id = $1"
	err := db.Client.Client.Get(&res, query, d.Id)
	if err == nil {
		res.LoadToDraft(d)
	}
	return err
}

func (d *Draft) FromLevelId() error {
	if d.LevelId == 0 {
		return errors.New("missing required field LevelId")
	}
	res := DbDraft{}
	query := "SELECT * FROM drafts WHERE level_id = $1"
	err := db.Client.Client.Get(&res, query, d.LevelId)
	if err == nil {
		res.LoadToDraft(d)
	}
	return err
}

func (d *Draft) Delete() error {
	if d.Id == 0 {
		return errors.New("missing required param id")
	}
	query := "DELETE FROM drafts WHERE id = $1;"
	_, err := db.Client.Client.Exec(query, d.Id)
	return err
}

func (d *Draft) SafeDelete() error {
	if d.Id == 0 || d.Uid == 0 {
		return errors.New("missing required param id, uid")
	}
	query := "DELETE FROM drafts WHERE id = $1 AND uid = $2;"
	_, err := db.Client.Client.Exec(query, d.Id, d.Uid)
	return err
}

func GetUserDrafts(uid int) ([]Draft, error) {
	query := "SELECT * FROM drafts WHERE uid = $1;"
	res := []DbDraft{}
	err := db.Client.Client.Select(&res, query, uid)
	arr := make([]Draft, 0, len(res))
	for _, v := range res {
		d := Draft{}
		v.LoadToDraft(&d)
		arr = append(arr, d)
	}
	return arr, err
}

// Returns the draft for the level.
// If no draft exists, returns a new
// draft. Can also be used to fork a
// level
func GetLevelDraft(levelId, uid int) (Draft, error) {
	d := Draft{LevelId: levelId}
	err := d.FromLevelId()
	if err != nil && err != sql.ErrNoRows {
		return Draft{}, err
	} else if err == nil && d.Uid == uid && d.Name != "" {
		return d, nil
	}
	level := Level{Id: levelId}
	err = level.Get()
	if err != nil {
		return Draft{}, err
	}
	name := level.Name
	if uid != level.Uid {
		name = fmt.Sprintf("Copia de %v", level.Name)
	}
	draft := Draft{
		Name:       name,
		Theme:      level.Theme,
		CourseData: level.CourseData,
		Uid:        uid,
		Car:        level.Car,
		Soundtrack: level.Soundtrack,
	}
	if uid == level.Uid {
		draft.LevelId = levelId
	}
	id, err := draft.Create()
	if err != nil {
		return Draft{}, err
	}
	draft.Id = id
	// We track here bc the route
	// doesn't know if a draft
	// was created
	event := tracking.Event{
		EventType: "draft_create",
		Uid:       uid,
		DraftId:   id,
	}
	event.Send()
	return draft, nil
}
