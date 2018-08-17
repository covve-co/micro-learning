package main

import (
	"fmt"
	"os"

	"github.com/covveco/micro-learning/db"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}

	if len(os.Args) != 2 {
		fmt.Println("need to specify 1 argument: up, down or force")
		return
	}

	if err := godotenv.Load(); err != nil {
		panic(".env file required")
	}
	log.Infoln("loaded .env file")

	db, err := db.New(os.Getenv("DB"))
	if err != nil {
		panic(err)
	}
	log.Infoln("opened connection to database")

	if err := db.Migrate(os.Args[1]); err != nil {
		panic(err)
	}
	log.Infoln("successfully migrated the database")
}
