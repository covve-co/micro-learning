package db

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // database driver
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
	"gitlab.com/covveco/special-needs/db/migrations"
)

type DB struct {
	db *sqlx.DB
}

// New creates a new database instance.
func New(url string) (*DB, error) {
	db := &DB{}

	sqldb, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	db.db = sqldb

	return db, nil
}

func (db *DB) Migrate(cmd string) error {
	as := bindata.Resource(migrations.AssetNames(), func(name string) ([]byte, error) {
		return migrations.Asset(name)
	})

	sd, err := bindata.WithInstance(as)
	if err != nil {
		return err
	}

	dd, err := postgres.WithInstance(db.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", sd, os.Getenv("DB_URL"), dd)
	if err != nil {
		return err
	}

	if cmd == "up" {
		if err := m.Up(); err != nil {
			return err
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil {
			return err
		}
	} else {
		v, err := strconv.Atoi(cmd)
		if err != nil {
			return err
		}

		if err := m.Force(v); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Seed(name string) error {
	var names []string

	if name == "all" {
		ns, err := filepath.Glob("db/seeds/*")
		if err != nil {
			return err
		}

		names = append(names, ns...)
	} else {
		names = append(names, "db/seeds/"+name+".sql")
	}

	for _, n := range names {
		b, err := ioutil.ReadFile(n)
		if err != nil {
			return err
		}

		_, err = db.db.Exec(string(b))
		if err != nil {
			return err
		}
	}

	return nil
}
