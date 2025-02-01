package db

import (
	"embed"
	"testing"
)

// go:embed migrations_test/*.sql
var testMigrationFS embed.FS

func TestTagDBInit(t *testing.T) {
	t.Error("test not written")
}

func TestTagDBClose(t *testing.T) {
	t.Error("test not written")
}
