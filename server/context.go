package server

import "github.com/wensboy/ss/db"

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
