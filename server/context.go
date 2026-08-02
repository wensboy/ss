package server

import (
	"github.com/wensboy/ss/config"
	"github.com/wensboy/ss/db"
)

type ServerContext struct {
	DBContext *db.SqlDBContext
}

func NewServerContext() *ServerContext {
	return &ServerContext{}
}

func (s *ServerContext) SetDBContext(dbContext *db.SqlDBContext) *ServerContext {
	s.DBContext = dbContext
	return s
}

// MountDBContext 挂载服务一致数据库上下文
func (s *ServerContext) MountDBContext() {
	// sql database context
	{
		dbType := config.MustLookup[string](
			config.GEnvSource("ss_db_type"),
			config.GConfigSource("db.type"),
			config.DefaultSource(db.DB_TYPE_SQLITE),
		)
		dbname := config.MustLookup[string](
			config.GEnvSource("ss_db_name"),
			config.GConfigSource("db.name"),
			config.DefaultSource("default"),
		)
		dsn := config.MustLookup[string](
			config.GEnvSource("ss_db_dsn"),
			config.GConfigSource("db.dsn"),
		)
		sqlDB, err := db.NewSqlDatabase(dbType, dbname, dsn)
		if err != nil {
			panic(err)
		}
		db.GetGSqlDBContext().Set(sqlDB)
		s.SetDBContext(db.GetGSqlDBContext())
	}
}

func (s *ServerContext) GetDBContext() *db.SqlDBContext {
	return s.DBContext
}
