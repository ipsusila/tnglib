package sqldb_test

import (
	"testing"

	"github.com/ipsusila/tnglib"
	"github.com/ipsusila/tnglib/sqldb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqldb(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	err = sqldb.RegisterDBX("testing", db)
	assert.NoError(t, err)

	err = tnglib.RunTengoScriptFile("../_testdata/sqldb.tengo", "sqldb", "context")
	assert.NoError(t, err)
}
