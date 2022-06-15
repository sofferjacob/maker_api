package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type client struct {
	Client *sqlx.DB
}

var Client client

func (c *client) Connect() {
	// dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Printf("❌ Error: could not connect to database %v @ %v:%v: %v\n", os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), err.Error())
		os.Exit(2)
	}
	c.Client = db
	// schemaFile, err := os.ReadFile("/home/jacobosoffer/projects/TC2005B/maker_api/sql/init.sql")
	schemaFile, err := os.ReadFile("sql/init.sql")
	if err != nil {
		fmt.Printf("❌ Error: Could not read schema file. %v\n", err.Error())
		os.Exit(2)
	}
	schema := string(schemaFile)
	_, err = db.Exec(schema)
	if err != nil {
		fmt.Printf("❌ Error: Could not load schema to db. %v\n", err.Error())
		os.Exit(2)
	}
}

func (c *client) Close() {
	c.Client.Close()
}
