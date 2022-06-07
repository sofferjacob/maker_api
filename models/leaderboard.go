package models

import (
	"time"

	"github.com/sofferjacob/maker_api/db"
)

type Leaderboard struct {
	LevelId   int       `db:"level_id" json:"levelId"`
	Uid       int       `db:"uid" json:"uid"`
	Time      int       `db:"time" json:"time"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Name      string    `db:"name" json:"userName"`
}

func GetLeaderboard(levelId int) ([]Leaderboard, error) {
	query := "SELECT * FROM leaderboard WHERE level_id = $1 LIMIT 10;"
	res := []Leaderboard{}
	err := db.Client.Client.Select(&res, query, levelId)
	return res, err
}
