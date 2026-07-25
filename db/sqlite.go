package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type SqliteDB struct {
	name string
	db   *xorm.Engine
}

func NewSqliteDB(name string, dsn string, driverOptions ...func(*sql.DB) error) (*SqliteDB, error) {
	sqliteDB := &SqliteDB{
		name: name,
	}
	var err error
	sqliteDB.db, err = xorm.NewEngine(sqliteDB.Driver(), dsn, driverOptions...)
	if err != nil {
		return nil, err
	}
	return sqliteDB, nil
}

func (s *SqliteDB) DB() *xorm.Engine {
	return s.db
}

func (s *SqliteDB) Name() string {
	return s.name
}

func (s *SqliteDB) DBName() string {
	return "sqlite"
}

func (s *SqliteDB) Driver() string {
	return "sqlite3"
}

func (s *SqliteDB) Sync(tables ...any) error {
	return s.db.Sync(tables...)
}

func (s *SqliteDB) Close() error {
	return s.db.Close()
}

func (s *SqliteDB) Ping() error {
	return s.db.Ping()
}
