package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/theflyingcodr/config"
)

type dbSetupFunc func(c *config.Db) (*sqlx.DB, error)
type dbSetups map[config.DbType]dbSetupFunc

// NewDbSetup will load the db setup functions into a lookup map
// ready for being called in main.go.
// nolint: revive // yep
func NewDbSetup() dbSetups {
	s := make(map[config.DbType]dbSetupFunc, 3)
	s[config.DBSqlite] = setupSqliteDB
	s[config.DBMySQL] = setupMySQLDB
	s[config.DBPostgres] = setupPostgresDB
	return s
}

// SetupDb can be used to setup a new database.
func (d dbSetups) SetupDb(cfg *config.Db) (*sqlx.DB, error) {
	fn, ok := d[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("db type %s not supported", cfg.Type)
	}
	return fn(cfg)
}
