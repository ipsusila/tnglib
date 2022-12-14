package tnglib

import (
	"context"
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

// standard context module
var ctxModule = map[string]tengo.Object{
	"background": &tengo.UserFunction{
		Name:  "background",
		Value: ctxFunc(context.Background),
	},
	"todo": &tengo.UserFunction{
		Name:  "todo",
		Value: ctxFunc(context.TODO),
	},
	"with_timeout": &tengo.UserFunction{
		Name:  "with_timeout",
		Value: ctxWithTimeout,
	},
	"with_cancel": &tengo.UserFunction{
		Name:  "with_cancel",
		Value: ctxWithCancel,
	},
	"with_deadline": &tengo.UserFunction{
		Name:  "with_deadline",
		Value: ctxWithDeadline,
	},
	"with_value": &tengo.UserFunction{
		Name:  "with_value",
		Value: ctxWithValue,
	},
}

// Context is context.Context wrapper which is accessible from tengo
type Context struct {
	tengo.ImmutableMap
	Ctx context.Context
}

// NewContext creates scriptable context.Context
func NewContext(ctx context.Context) *Context {
	return &Context{
		Ctx: ctx,
		ImmutableMap: tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"value": &tengo.UserFunction{
					Name:  "value",
					Value: FuncAIRI(ctx.Value),
				},
				"err": &tengo.UserFunction{
					Name:  "err",
					Value: stdlib.FuncARE(ctx.Err),
				},
				"deadline": &tengo.UserFunction{
					Name:  "deadline",
					Value: FuncATBR(ctx.Deadline),
				},
			},
		},
	}
}

// TypeName return context
func (c *Context) TypeName() string {
	return "context"
}

// String return string representation of the context
func (c *Context) String() string {
	return fmt.Sprintf("<context>:%v, map: %s", c.Ctx, c.ImmutableMap.String())
}

// Copy returns a copy of the type.
func (c *Context) Copy() tengo.Object {
	return &Context{Ctx: c.Ctx, ImmutableMap: c.ImmutableMap}
}

// IsFalsy returns true if the value of the type is falsy.
func (c *Context) IsFalsy() bool {
	return c.Ctx == nil
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (c *Context) Equals(x tengo.Object) bool {
	if x == nil || c == x {
		return c == x
	}
	v, ok := x.(*Context)
	if !ok {
		return false
	}
	return v.Ctx == c.Ctx
}

func ctxFunc(fn func() context.Context) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		return NewContext(fn()), nil
	}
}

func ctxWithCancel(args ...tengo.Object) (tengo.Object, error) {
	parent, err := ArgToContext(args...)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parent.Ctx)
	return makeCtxReturn(ctx, cancel), nil
}

func ctxWithDeadline(args ...tengo.Object) (tengo.Object, error) {
	parent, err := ArgIToContext(0, args...)
	if err != nil {
		return nil, err
	}
	deadline, err := ArgIToTime(1, args...)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithDeadline(parent.Ctx, deadline)
	return makeCtxReturn(ctx, cancel), nil
}

func ctxWithTimeout(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	// get context
	parent, err := ArgIToContext(0, args...)
	if err != nil {
		return nil, err
	}

	// get timeout value
	timeout, cverr, argerr := ArgIToDuration(1, args...)
	if argerr != nil {
		return nil, err
	} else if cverr != nil {
		return WrapError(cverr), nil
	}

	ctx, cancel := context.WithTimeout(parent.Ctx, timeout)
	return makeCtxReturn(ctx, cancel), nil
}

func ctxWithValue(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 3 {
		return nil, tengo.ErrWrongNumArguments
	}
	parent, err := ArgIToContext(0, args...)
	if err != nil {
		return nil, err
	}
	key := tengo.ToInterface(args[1])
	val := tengo.ToInterface(args[2])
	ctx := context.WithValue(parent.Ctx, key, val)
	return NewContext(ctx), nil
}

func makeCtxReturn(ctx context.Context, cancel context.CancelFunc) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"ctx": NewContext(ctx),
			"cancel": &tengo.UserFunction{
				Name:  "cancel",
				Value: stdlib.FuncAR(cancel),
			},
		},
	}
}
