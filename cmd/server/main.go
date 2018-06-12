package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/covveco/special-needs/db"
	"gitlab.com/covveco/special-needs/handler"
	"gitlab.com/covveco/special-needs/view"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		DisableColors: true,
	}

	if err := godotenv.Load(); err != nil {
		panic(".env file required")
	}
	log.Infoln("loaded .env file")

	var out io.Writer = os.Stdout

	fname := os.Getenv("LOG_FILE")
	if fname != "" {
		f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic("couldn't open log file")
		}
		defer f.Close()

		// Tee the output
		out = io.MultiWriter(os.Stdout, f)
	}

	log.Out = out

	if err := view.Parse(); err != nil {
		panic(errors.Wrap(err, "view: failed to parse"))
	}
	log.Infoln("parsed templates and layouts")

	db, err := db.New(os.Getenv("DB"))
	if err != nil {
		panic(errors.Wrap(err, "db: failed to connect"))
	}
	log.Infoln("opened connection to database")

	handler := handler.Handler{
		DB:            db,
		Log:           log,
		AuthSecret:    os.Getenv("AUTH_SECRET"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}

	log.Infoln(fmt.Sprintf("server started listening on %s...", os.Getenv("PORT")))
	if err := http.ListenAndServe(os.Getenv("PORT"), handler.Routes()); err != nil {
		log.Infoln("shutting down server")
		panic(err)
	}
}
