package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"xorm.io/xorm"
)

type MariaDB struct {
	name string
	db   *xorm.Engine
}

func NewMariaDB(name string, dsn string, driverOptions ...func(*sql.DB) error) (*MariaDB, error) {
	mariaDB := &MariaDB{
		name: name,
	}
	var err error
	mariaDB.db, err = xorm.NewEngine(mariaDB.Driver(), dsn, driverOptions...)
	if err != nil {
		return nil, err
	}
	return mariaDB, nil
}

func (m *MariaDB) DB() *xorm.Engine {
	return m.db
}

func (m *MariaDB) Name() string {
	return m.name
}

func (m *MariaDB) DBName() string {
	return "mariadb"
}

func (m *MariaDB) Driver() string {
	return "mysql"
}

func (m *MariaDB) Sync(tables ...any) error {
	return m.db.Sync(tables...)
}

func (m *MariaDB) Close() error {
	return m.db.Close()
}

func (m *MariaDB) Ping() error {
	return m.db.Ping()
}
