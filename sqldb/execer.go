package sqldb

import (
	"context"
	"database/sql"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeExecer(ex sqlx.ExecerContext) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"exec": &tengo.UserFunction{
				Name:  "exec",
				Value: exec(ex.ExecContext),
			},
		},
	}
}

func execerFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := tnglib.ArgToString(args...)
		if err != nil {
			return nil, err
		}
		dbx, err := sqldbRegistry.Entry(s)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeExecer(dbx), nil
	}
}

func exec(fn func(context.Context, string, ...interface{}) (sql.Result, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		ctx, err := tnglib.ArgIToContext(0, args...)
		if err != nil {
			return nil, err
		}
		query, err := tnglib.ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}

		qargs := []any{}
		for i := 2; i < len(args); i++ {
			qargs = append(qargs, tengo.ToInterface(args[i]))
		}

		res, err := fn(ctx.Value, query, qargs...)
		if err != nil {
			return tnglib.WrapError(err), nil
		}

		m := make(map[string]tengo.Object)
		if id, err := res.LastInsertId(); err == nil {
			m["last_insert_id"] = &tengo.Int{Value: id}
		}
		if nr, err := res.RowsAffected(); err == nil {
			m["rows_affected"] = &tengo.Int{Value: nr}
		}
		return &tengo.Map{Value: m}, nil
	}
}
