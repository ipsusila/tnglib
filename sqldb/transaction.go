package sqldb

import (
	"context"
	"database/sql"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeTransaction(tx *sqlx.Tx) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"driver_name": &tengo.UserFunction{
				Name:  "driver_name",
				Value: stdlib.FuncARS(tx.DriverName),
			},
			"one": &tengo.UserFunction{
				Name:  "one",
				Value: queryRowx(tx.QueryRowxContext),
			},
			"many": &tengo.UserFunction{
				Name:  "many",
				Value: queryRowsx(tx.QueryxContext),
			},
			"exec": &tengo.UserFunction{
				Name:  "exec",
				Value: exec(tx.ExecContext),
			},
			"rebind": &tengo.UserFunction{
				Name:  "rebind",
				Value: stdlib.FuncASRS(tx.Rebind),
			},
			"bind_named": &tengo.UserFunction{
				Name:  "bind_named",
				Value: tnglib.FuncASARSAE(tx.BindNamed),
			},
			"commit": &tengo.UserFunction{
				Name:  "commit",
				Value: stdlib.FuncARE(tx.Commit),
			},
			"rollback": &tengo.UserFunction{
				Name:  "rollback",
				Value: stdlib.FuncARE(tx.Rollback),
			},
		},
	}
}

func transactionFunc(fn0 func() (*sqlx.Tx, error), fn2 func(context.Context, *sql.TxOptions) (*sqlx.Tx, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 0 {
			if fn0 == nil {
				return nil, tengo.ErrWrongNumArguments
			}

			tx, err := fn0()
			if err != nil {
				return tnglib.WrapError(err), nil
			}
			return makeTransaction(tx), nil
		}

		// get context
		ctx, err := tnglib.ArgIToContext(0, args...)
		if err != nil {
			return nil, err
		}

		// get transaction options
		// func(ctx, isolation_level, [read_only])
		var opts *sql.TxOptions
		if len(args) >= 2 {
			vi, err := tnglib.ArgIToInt(1, args...)
			if err != nil {
				return nil, err
			}
			opts = &sql.TxOptions{
				Isolation: sql.IsolationLevel(vi),
			}
		}
		if len(args) == 3 {
			vb, err := tnglib.ArgIToBool(2, args...)
			if err != nil {
				return nil, err
			}
			opts.ReadOnly = vb
		}

		tx, err := fn2(ctx.Value, opts)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeTransaction(tx), nil
	}
}
