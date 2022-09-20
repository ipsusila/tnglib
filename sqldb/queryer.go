package sqldb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	"github.com/jmoiron/sqlx"
)

func makeQueryer(qr sqlx.QueryerContext) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"one": &tengo.UserFunction{
				Name:  "one",
				Value: queryRowx(qr.QueryRowxContext),
			},
			"many": &tengo.UserFunction{
				Name:  "many",
				Value: queryRowsx(qr.QueryxContext),
			},
		},
	}
}

func queryerFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := tnglib.ArgToString(args...)
		if err != nil {
			return nil, err
		}
		dbx, err := sqldbRegistry.Entry(s)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeQueryer(dbx), nil
	}
}

func queryRowx(fn func(context.Context, string, ...interface{}) *sqlx.Row) tengo.CallableFunc {
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

		row := fn(ctx.Ctx, query, qargs...)
		result := make(map[string]interface{})
		if err := row.MapScan(result); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return tengo.FalseValue, nil
			} else {
				return tnglib.WrapError(err), nil
			}
		}
		return tnglib.MapToObject(result)
	}
}

func queryRowsx(fn func(context.Context, string, ...interface{}) (*sqlx.Rows, error)) tengo.CallableFunc {
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

		rows, err := fn(ctx.Ctx, query, qargs...)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		defer rows.Close()

		results := []tengo.Object{}
		for rows.Next() {
			result := make(map[string]interface{})
			if err := rows.MapScan(result); err != nil {
				return tnglib.WrapError(err), nil
			}
			obj, err := tnglib.MapToObject(result)
			if err != nil {
				return tnglib.WrapError(err), nil
			}
			results = append(results, obj)
		}

		// return result as array of items
		return &tengo.Array{Value: results}, nil
	}
}
