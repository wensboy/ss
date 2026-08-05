package server

import (
	"github.com/wensboy/ss/config"
	"github.com/wensboy/ss/db"
	"github.com/wensboy/ss/err"
	"github.com/wensboy/ss/log"
)

type ServerContext struct {
	DBContext *db.SqlDBContext
	Logger    *log.MutateLogger
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
		sqlDB, perr := db.NewSqlDatabase(dbType, dbname, dsn)
		if perr != nil {
			panic(err.GetGErrHub().GetOrSet(err.NewNormalErrCode(ErrCodeRequiredOption), err.NewErr("", err.NewNormalErrCode(ErrCodeRequiredOption).Code(), "database init failed").Wrap(perr)))
		}
		db.GetGSqlDBContext().Set(sqlDB)
		s.SetDBContext(db.GetGSqlDBContext())
	}
}

func (s *ServerContext) GetDBContext() *db.SqlDBContext {
	return s.DBContext
}

func (s *ServerContext) SetLogger(logger *log.MutateLogger) *ServerContext {
	s.Logger = logger
	return s
}

func (s *ServerContext) MountLogger() {
	logLevel := config.MustLookup[string](
		config.GEnvSource("ss_log_level"),
		config.GConfigSource("server.log.level"),
		config.DefaultSource("info"),
	)
	logEnv := config.MustLookup[string](
		config.GEnvSource("ss_log_env"),
		config.GConfigSource("server.log.env"),
		config.DefaultSource("development"),
	)
	logger := log.NewMutateLogger(
		log.WithLogLevel(s.levelCast(logLevel)),
		log.WithLogEnv(s.envCast(logEnv)),
	)
	s.SetLogger(logger)
}

func (s *ServerContext) GetLogger() *log.MutateLogger {
	return s.Logger
}

func (s *ServerContext) levelCast(lv string) log.Level {
	switch lv {
	case "debug":
		return log.LevelDebug
	case "info":
		return log.LevelInfo
	case "warn":
		return log.LevelWarn
	case "error":
		return log.LevelError
	case "panic":
		return log.LevelPanic
	case "fatal":
		return log.LevelFatal
	default:
		return log.LevelInfo
	}
}

func (s *ServerContext) envCast(env string) log.LogEnv {
	switch env {
	case "debug":
		return log.LogEnvDebug
	case "development":
		return log.LogEnvDevelopment
	case "example":
		return log.LogEnvExample
	case "production":
		return log.LogEnvProduction
	default:
		return log.LogEnvDevelopment
	}
}
