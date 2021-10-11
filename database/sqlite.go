//go:build sqlite
// Package database defines db level setup.
// nolint: dupl // fine to duplicate for now
package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/theflyingcodr/config"

	// used to import the sqlite drivers.
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func setupSqliteDB(c *config.Db) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", c.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup database")
	}
	if !c.Migrate {
		log.Println("migrate database set to false, skipping migration")
		return db, nil
	}
	log.Println("migrating database")
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("creating sqlite3 db driver failed %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", c.SchemaPath), "sqlite3",
		driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}

	}
	log.Println("migrating database completed")
	return db, nil
}
