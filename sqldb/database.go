package sqldb

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeDatabase(db *sqlx.DB) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"driver_name": &tengo.UserFunction{
				Name:  "driver_name",
				Value: stdlib.FuncARS(db.DriverName),
			},
			"one": &tengo.UserFunction{
				Name:  "one",
				Value: queryRowx(db.QueryRowxContext),
			},
			"many": &tengo.UserFunction{
				Name:  "many",
				Value: queryRowsx(db.QueryxContext),
			},
			"exec": &tengo.UserFunction{
				Name:  "exec",
				Value: exec(db.ExecContext),
			},
			"rebind": &tengo.UserFunction{
				Name:  "rebind",
				Value: stdlib.FuncASRS(db.Rebind),
			},
			"bind_named": &tengo.UserFunction{
				Name:  "bind_named",
				Value: tnglib.FuncASARSAE(db.BindNamed),
			},
			"transaction": &tengo.UserFunction{
				Name:  "transaction",
				Value: transactionFunc(db),
			},
			"ping": &tengo.UserFunction{
				Name:  "ping",
				Value: tnglib.FuncACRE(db.PingContext),
			},
		},
	}
}

func databaseFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := tnglib.ArgToString(args...)
		if err != nil {
			return nil, err
		}
		dbx, err := sqldbRegistry.Entry(s)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeDatabase(dbx), nil
	}
}
