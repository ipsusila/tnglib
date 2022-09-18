package sqldb

import (
	"database/sql"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/registry"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

// Name for this module
const (
	Name = "sqldb"
)

// module definition
var (
	sqldbRegistry = registry.NewSyncMapRegistry[string, *sqlx.DB]()

	// Module registered here
	sqldbModule = map[string]tengo.Object{
		// place holder
		"unknown":  &tengo.Int{Value: sqlx.UNKNOWN},
		"question": &tengo.Int{Value: sqlx.QUESTION},
		"dollar":   &tengo.Int{Value: sqlx.DOLLAR},
		"named":    &tengo.Int{Value: sqlx.NAMED},
		"at":       &tengo.Int{Value: sqlx.AT},
		// isolation level
		"level_default":          &tengo.Int{Value: int64(sql.LevelDefault)},
		"level_read_uncommitted": &tengo.Int{Value: int64(sql.LevelReadUncommitted)},
		"level_read_committed":   &tengo.Int{Value: int64(sql.LevelReadCommitted)},
		"level_write_committed":  &tengo.Int{Value: int64(sql.LevelWriteCommitted)},
		"level_repeatable_read":  &tengo.Int{Value: int64(sql.LevelRepeatableRead)},
		"level_snapshot":         &tengo.Int{Value: int64(sql.LevelSnapshot)},
		"level_serializable":     &tengo.Int{Value: int64(sql.LevelSerializable)},
		"level_linearizable":     &tengo.Int{Value: int64(sql.LevelLinearizable)},

		// database(string) => db
		"database": &tengo.UserFunction{
			Name:  "database",
			Value: databaseFunc(),
		},

		// queryer(string) => queryer
		"queryer": &tengo.UserFunction{
			Name:  "queryer",
			Value: queryerFunc(),
		},

		// execer(string) => execer
		"execer": &tengo.UserFunction{
			Name:  "execer",
			Value: execerFunc(),
		},

		// bind_named(int, string, interface{}) => ([]{string, []interface{}}, error)
		"bind_named": &tengo.UserFunction{
			Name:  "bind_named",
			Value: tnglib.FuncAISARSAE(sqlx.BindNamed),
		},
		// bind_named(int, string) => string
		"rebind": &tengo.UserFunction{
			Name:  "rebind",
			Value: tnglib.FuncAISRS(sqlx.Rebind),
		},

		// expand_in(string, ...any) => ([]{string, []interface{}}, error)
		"expand_in": &tengo.UserFunction{
			Name:  "expand_in",
			Value: tnglib.FuncASVRSAE(sqlx.In),
		},
	}
)

func init() {
	// register module
	tnglib.RegisterModule(Name, sqldbModule)
}

// RegisterDBX register sqlx.DB instance
func RegisterDBX(name string, dbx *sqlx.DB) error {
	return sqldbRegistry.Register(name, dbx)
}

// RegisterDB register existing sql.DB instance. Supported driver names:
//   - With $ params: "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach",
//   - With ? params: "mysql", "sqlite3", "nrmysql", "nrsqlite3",
//   - with :name params: "oci8", "ora", "goracle", "godror",
//   - with @ params: "sqlserver",
func RegisterDB(name, driverName string, db *sql.DB) error {
	dbx := sqlx.NewDb(db, driverName)
	return RegisterDBX(name, dbx)
}

// ReplaceDBX with new instance
func ReplaceDBX(name string, dbx *sqlx.DB) (*sqlx.DB, error) {
	return sqldbRegistry.Replace(name, dbx)
}

// ReplaceDB with new instance
func ReplaceDB(name, driverName string, db *sql.DB) (*sqlx.DB, error) {
	dbx := sqlx.NewDb(db, driverName)
	return ReplaceDBX(name, dbx)
}

// Module stores information for this module name
func Module() (string, tengo.Importable) {
	return Name, &tengo.BuiltinModule{Attrs: sqldbModule}
}
