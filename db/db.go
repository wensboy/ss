package db

import (
	"sync"

	"xorm.io/xorm"
)

var (
	_global_db_context *SqlDBContext
)

type SqlTable interface {
	TableName() string
}

type SqlDatabase interface {
	Name() string
	DBName() string
	Driver() string
	DB() *xorm.Engine
	Sync(...any) error
	Close() error
	Ping() error
}

func init() {
	_global_db_context = NewSqlDBContext()
}

func To[T SqlDatabase](db SqlDatabase) (T, bool) {
	var zero T
	if db == nil {
		return zero, false
	}
	if t, ok := db.(T); ok {
		return t, true
	}
	return zero, false
}

func MustTo[T SqlDatabase](db SqlDatabase) T {
	t, ok := To[T](db)
	if !ok {
		panic("db type assertion failed")
	}
	return t
}

func GetGSqlDBContext() *SqlDBContext {
	return _global_db_context
}

func SetGSqlDBContext(c *SqlDBContext) {
	_global_db_context = c
}

type SqlDBContext struct {
	dbs map[string]SqlDatabase
	mu  sync.RWMutex
}

func NewSqlDBContext() *SqlDBContext {
	return &SqlDBContext{
		dbs: make(map[string]SqlDatabase),
	}
}

func (c *SqlDBContext) Set(db SqlDatabase) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dbs[db.DBName()+"::"+db.Name()] = db
}

func (c *SqlDBContext) Get(dbname, name string) (SqlDatabase, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	db, ok := c.dbs[dbname+"::"+name]
	return db, ok
}

func (c *SqlDBContext) MustGet(dbname, name string) SqlDatabase {
	db, ok := c.Get(dbname, name)
	if !ok {
		panic("db not found: " + dbname + "::" + name)
	}
	return db
}

func (c *SqlDBContext) Del(dbname, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.dbs, dbname+"::"+name)
}

func (c *SqlDBContext) Ping() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, db := range c.dbs {
		if err := db.Ping(); err != nil {
			return err
		}
	}
	return nil
}
