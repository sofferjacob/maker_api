package tracking

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/sofferjacob/maker_api/db"
)

type Event struct {
	Id        int                    `db:"id" json:"id"`
	EventType string                 `db:"event_type" json:"eventType" binding:"required"`
	LevelId   int                    `db:"level_id" json:"levelId"`
	Timestamp time.Time              `db:"timestamp" json:"timestamp"`
	Uid       int                    `db:"uid" json:"uid"`
	Time      int                    `db:"time" json:"time"`
	DraftId   int                    `db:"draft_id" json:"draft_id"`
	Body      map[string]interface{} `db:"body" json:"body"`
	State     string                 `db:"state" json:"state"`
}

func (e *Event) Send() error {
	if e.EventType == "" {
		return errors.New("field EventType must not be null")
	}
	qb := db.Insert("events").Set("event_type", e.EventType)
	if e.LevelId != 0 {
		qb = qb.Set("level_id", e.LevelId)
	}
	if e.Uid != 0 {
		qb = qb.Set("uid", e.Uid)
	}
	if e.Time != 0 {
		qb = qb.Set("time", e.Time)
	}
	if e.DraftId != 0 {
		qb = qb.Set("draft_id", e.DraftId)
	}
	if e.Body != nil || len(e.Body) > 0 {
		body, err := json.Marshal(e.Body)
		if err != nil {
			return fmt.Errorf("could not parse body: %v", err.Error())
		}
		qb = qb.Set("body", body)
	}
	if e.State != "" {
		qb = qb.Set("state", e.State)
	}
	query, args := qb.Query()
	_, err := db.Client.Client.Exec(query, args...)
	return err
}
