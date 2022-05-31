package conf

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var params []string = []string{
	"PORT",
	"AUTH_KEY",
	"DB_PORT",
	"DB_HOST",
	"DB_USER",
	"DB_PASSWORD",
	"DB_NAME",
}

func validate() {
	for _, v := range params {
		if os.Getenv(v) == "" {
			fmt.Printf("Error: Required config parameter %v not found, terminating\n", v)
			os.Exit(2)
		}
	}
}

func Load() {
	godotenv.Load()
	validate()
}
