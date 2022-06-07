package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sofferjacob/maker_api/db"
)

type Collection struct {
	Id          int          `json:"id" db:"id"`
	Name        string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	Uid         int          `db:"uid" json:"uid"`
	Created     time.Time    `db:"created" json:"created"`
	Updated     sql.NullTime `db:"updated" json:"updated"`
	Ts          string       `db:"ts"`
}

type CollectionData struct {
	Collection
	UserName string `db:"user_name" json:"userName"`
}

func (c *Collection) Create() (int, error) {
	if c.Name == "" || c.Uid == 0 {
		return 0, errors.New("missing required fields (uid, name)")
	}
	var id int
	query := "INSERT INTO collection (uid, name, description) VALUES ($1, $2, $3) RETURNING id;"
	err := db.Client.Client.Get(&id, query, c.Uid, c.Name, c.Description)
	return id, err
}

func GetCollection(id int) (CollectionData, error) {
	query := "SELECT c.*, u.name user_name FROM collection c INNER JOIN users u ON uid = u.id WHERE c.id = $1;"
	res := CollectionData{}
	err := db.Client.Client.Get(&res, query, id)
	return res, err
}

func GetUserCollections(uid int) ([]CollectionData, error) {
	query := "SELECT c.*, u.name user_name FROM collection c INNER JOIN users u ON uid = u.id WHERE c.uid = $1;"
	res := []CollectionData{}
	err := db.Client.Client.Select(&res, query, uid)
	return res, err
}

func (c *Collection) Delete() error {
	if c.Id == 0 || c.Uid == 0 {
		return errors.New("missing required fields id and uid")
	}
	query := "DELETE FROM collection WHERE id = $1 AND uid = $2;"
	_, err := db.Client.Client.Exec(query, c.Id, c.Uid)
	return err
}

func (c *Collection) Update() error {
	if c.Id == 0 || c.Uid == 0 {
		return errors.New("missing required fields Id or Uid")
	}
	if c.Name == "" && c.Description == "" {
		return nil
	}
	i := 1
	args := []interface{}{}
	addComma := false
	query := "UPDATE collection SET "
	if c.Name != "" {
		query += fmt.Sprintf("name = $%v", i)
		i++
		addComma = true
		args = append(args, c.Name)
	}
	if c.Description != "" {
		if addComma {
			query += ","
		}
		query += fmt.Sprintf(" description = $%v", i)
		i++
		args = append(args, c.Description)
	}
	query += fmt.Sprintf(" WHERE id=$%v AND uid=$%v", i, i+1)
	args = append(args, c.Id, c.Uid)
	_, err := db.Client.Client.Exec(query, args...)
	return err
}

func QueryCollectionsFTS(query string) ([]Collection, error) {
	res := []Collection{}
	err := db.Client.Client.Select(&res, "SELECT * FROM query_gin(null::collection, $1);", query)
	return res, err
}

func IsOwnCollection(collectionId, uid int) (bool, error) {
	query := "SELECT uid FROM collection WHERE id = $1;"
	var cuid int
	err := db.Client.Client.Get(&cuid, query, collectionId)
	return cuid == uid, err
}

func TrendingCollections() ([]Collection, error) {
	query := "SELECT id, name, description, uid, created, updated FROM trending_collections;"
	res := []Collection{}
	err := db.Client.Client.Select(&res, query)
	return res, err
}
