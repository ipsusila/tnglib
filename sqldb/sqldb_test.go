package sqldb_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/ipsusila/tnglib/sqldb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func execTest(t *testing.T, driverName, dsn string, scriptFile string) {
	if scriptFile == "" {
		scriptFile = "../_testdata/sqldb.tengo"
	}
	db, err := sqlx.Open(driverName, dsn)
	assert.NoError(t, err)
	defer db.Close()

	err = sqldb.RegisterDBX("testing", db)
	assert.NoError(t, err)

	err = script.RunFile(scriptFile, "fmt", "sqldb", "context")
	assert.NoError(t, err)
}

func TestGoSqlite(t *testing.T) {
	execTest(t, "sqlite3", ":memory:", "")
}
func TestPureGoSqlite(t *testing.T) {
	driverName := "sqlite"
	sqlx.BindDriver(driverName, sqlx.QUESTION)

	execTest(t, driverName, ":memory:", "")
}

func TestOpenSqlite(t *testing.T) {
	driverName := "sqlite"
	sqlx.BindDriver(driverName, sqlx.QUESTION)

	execTest(t, driverName, ":memory:", "../_testdata/sqldb_open.tengo")
}
