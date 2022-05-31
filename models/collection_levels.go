package models

import (
	"errors"

	"github.com/sofferjacob/maker_api/db"
)

type CollectionLevels struct {
	Id           int `db:"id"`
	CollectionId int `db:"collection_id"`
	LevelId      int `db:"level_id"`
}

func (c *CollectionLevels) Create() error {
	if c.CollectionId == 0 || c.LevelId == 0 {
		return errors.New("missing required fields")
	}
	query, args := db.Insert("collection_levels").
		Set("collection_id", c.CollectionId).
		Set("level_id", c.LevelId).Query()
	_, err := db.Client.Client.Exec(query, args...)
	return err
}

func (c *CollectionLevels) Delete() error {
	if c.CollectionId == 0 || c.LevelId == 0 {
		return errors.New("missing required fields")
	}
	query := "DELETE FROM collection_levels WHERE collection_id = $1 AND level_id = $2;"
	_, err := db.Client.Client.Exec(query, c.CollectionId, c.LevelId)
	return err
}

func GetCollectionLevels(id int) ([]Level, error) {
	query := "SELECT l.* FROM collection_levels c INNER JOIN levels l ON c.level_id = l.id WHERE c.collection_id = $1;"
	res := []DBLevel{}
	err := db.Client.Client.Select(&res, query, id)
	if err != nil {
		return nil, err
	}
	levels := make([]Level, 0, len(res))
	for _, v := range res {
		l := Level{}
		v.ToLevel(&l)
		levels = append(levels, l)
	}
	return levels, nil
}
