package db

import (
	"testing"

	"github.com/spf13/cast"
)

type TestTable struct {
	ID   int64
	Name string
}

func (t *TestTable) TableName() string {
	return "test_table"
}

func Test_CastSlice(t *testing.T) {
	a := []SqlTable{&TestTable{Name: "Alice"}, &TestTable{Name: "Bob"}, &TestTable{Name: "Charlie"}}
	b := cast.ToSlice(a)
	t.Logf("b: %v", b)
	if len(b) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(b))
	}
}

func Test_SqlDBContext(t *testing.T) {
	dbContext := GetGSqlDBContext()
	if dbContext == nil {
		t.Error("Expected non-nil db context")
	}
	sqliteDB, err := NewSqliteDB("test_sqlite", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatalf("Failed to create SqliteDB: %v", err)
	}
	dbContext.Set(sqliteDB)
	retrievedDB, ok := dbContext.Get("sqlite", "test_sqlite")
	if !ok {
		t.Error("Expected to retrieve the added database")
	}
	if retrievedDB.Name() != "test_sqlite" || retrievedDB.DBName() != "sqlite" {
		t.Errorf("Expected database name 'test_sqlite' and DB name 'sqlite', got '%s' and '%s'", retrievedDB.Name(), retrievedDB.DBName())
	}
	dbContext.Del("sqlite", "test_sqlite")
	_, ok = dbContext.Get("sqlite", "test_sqlite")
	if ok {
		t.Error("Expected database to be removed")
	}
}
