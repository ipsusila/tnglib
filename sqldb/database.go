package sqldb

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeDatabase(db *sqlx.DB, closable bool) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"driver_name": &tengo.UserFunction{
				Name:  "driver_name",
				Value: stdlib.FuncARS(db.DriverName),
			},
			"conn": &tengo.UserFunction{
				Name:  "conn",
				Value: connFunc(db),
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
				Value: transactionFunc(db.Beginx, db.BeginTxx),
			},
			"ping": &tengo.UserFunction{
				Name:  "ping",
				Value: tnglib.FuncACRE(db.PingContext),
			},
			"close": &tengo.UserFunction{
				Name:  "close",
				Value: closeableFunc(db, closable),
			},
		},
	}
}

func databaseFunc(closable bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := tnglib.ArgToString(args...)
		if err != nil {
			return nil, err
		}
		dbx, err := sqldbRegistry.Entry(s)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeDatabase(dbx, closable), nil
	}
}
func closeableFunc(db *sqlx.DB, closable bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		if !closable {
			// NOP for non closable db
			return tengo.TrueValue, nil
		}
		return tnglib.WrapError(db.Close()), nil
	}
}
func connectFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}

		ctx, err := tnglib.ArgIToContext(0, args...)
		if err != nil {
			return nil, err
		}

		drvName, err := tnglib.ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}
		dsn, err := tnglib.ArgIToString(2, args...)
		if err != nil {
			return nil, err
		}

		db, err := sqlx.ConnectContext(ctx.Value, drvName, dsn)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeDatabase(db, true), nil
	}
}

func openFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		drvName, err := tnglib.ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		dsn, err := tnglib.ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}

		db, err := sqlx.Open(drvName, dsn)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeDatabase(db, true), nil
	}
}
