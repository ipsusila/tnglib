package sqldb

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeConn(con *sqlx.Conn) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"one": &tengo.UserFunction{
				Name:  "one",
				Value: queryRowx(con.QueryRowxContext),
			},
			"many": &tengo.UserFunction{
				Name:  "many",
				Value: queryRowsx(con.QueryxContext),
			},
			"exec": &tengo.UserFunction{
				Name:  "exec",
				Value: exec(con.ExecContext),
			},
			"rebind": &tengo.UserFunction{
				Name:  "rebind",
				Value: stdlib.FuncASRS(con.Rebind),
			},
			"transaction": &tengo.UserFunction{
				Name:  "transaction",
				Value: transactionFunc(nil, con.BeginTxx),
			},
			"ping": &tengo.UserFunction{
				Name:  "ping",
				Value: tnglib.FuncACRE(con.PingContext),
			},
			"close": &tengo.UserFunction{
				Name:  "close",
				Value: stdlib.FuncARE(con.Close),
			},
		},
	}
}

func connFunc(db *sqlx.DB) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		cv, err := tnglib.ArgToContext(args...)
		if err != nil {
			return nil, err
		}
		con, err := db.Connx(cv.Ctx)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeConn(con), nil
	}
}
