package db

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/wensboy/ss/err"
	"xorm.io/xorm"
)

const (
	_default_dbcontext_sep = "::"
)

var (
	_global_db_context *SqlDBContext

	ErrCodeUnmatchedType = [4]int{err.Unwrap, err.PackageDB, 1, 1}
	ErrCodeDBNotFound    = [4]int{err.Unwrap, err.PackageDB, 1, 2}
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
	_global_db_context = NewSqlDBContext(_default_dbcontext_sep)
}

func NewSqlDatabase(dbname, name, dsn string, driverOptions ...func(*sql.DB) error) (SqlDatabase, error) {
	switch dbname {
	case DB_TYPE_SQLITE:
		return NewSqliteDB(name, dsn, driverOptions...)
	case DB_TYPE_MARIADB:
		return NewMariaDB(name, dsn, driverOptions...)
	default:
		return nil, errors.New("unsupported database type")
	}
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
		panic(err.GetGErrHub().GetOrSet(err.NewNormalErrCode(ErrCodeUnmatchedType), err.NewErr("", err.NewNormalErrCode(ErrCodeUnmatchedType).Code(), "unmatched db type")))
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
	sep string
	mu  sync.RWMutex
}

func NewSqlDBContext(sep string) *SqlDBContext {
	if sep == "" {
		sep = _default_dbcontext_sep
	}
	return &SqlDBContext{
		dbs: make(map[string]SqlDatabase),
		sep: sep,
	}
}

func (c *SqlDBContext) Set(db SqlDatabase) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dbs[db.DBName()+c.sep+db.Name()] = db
}

func (c *SqlDBContext) Get(dbname, name string) (SqlDatabase, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	db, ok := c.dbs[dbname+c.sep+name]
	return db, ok
}

func (c *SqlDBContext) MustGet(dbname, name string) SqlDatabase {
	db, ok := c.Get(dbname, name)
	if !ok {
		panic(err.GetGErrHub().GetOrSet(err.NewNormalErrCode(ErrCodeDBNotFound), err.NewErr("", err.NewNormalErrCode(ErrCodeDBNotFound).Code(), "db(%s) not found", dbname+c.sep+name)))
	}
	return db
}

func (c *SqlDBContext) Del(dbname, name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.dbs, dbname+c.sep+name)
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
