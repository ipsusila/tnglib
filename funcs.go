package tnglib

import (
	"io"
	"strings"

	"github.com/d5/tengo/v2"
)

// FuncWIE convert any function(io.Writer, interface{}) error to tengo.CallableFunc
func FuncWIE(fn func(w io.Writer, data any) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}

		// get writer interface
		var wr io.Writer
		if w, ok := args[0].(*Writer); ok {
			wr = w.Value
		} else {
			wr, ok = NewIoFunc(mWriter, args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "io.Writer",
					Found:    args[0].TypeName(),
				}
			}
		}
		da := tengo.ToInterface(args[1])
		return WrapError(fn(wr, da)), nil
	}
}

// FuncWISE convert any function(io.Writer, interface{}) error to tengo.CallableFunc
func FuncWISE(fn func(w io.Writer, data any) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}

		var sb strings.Builder
		da := tengo.ToInterface(args[0])
		if err := fn(&sb, da); err != nil {
			return WrapError(err), nil
		}
		return &tengo.String{Value: sb.String()}, nil
	}
}

// FuncARBs transform a function of 'func() []byte' signature into
// CallableFunc type.
func FuncARBs(fn func() []byte) tengo.CallableFunc {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		return &tengo.Bytes{Value: fn()}, nil
	}
}

// FuncBI transform a function of 'func([]byte, int)' signature into
// CallableFunc type.
func FuncBI(fn func([]byte, int)) tengo.CallableFunc {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		buf, err := ArgIToByteSlice(0, args...)
		if err != nil {
			return nil, err
		}
		n, err := ArgIToInt(1, args...)
		if err != nil {
			return nil, err
		}
		fn(buf, n)

		return tengo.TrueValue, nil
	}
}

// FuncE transform a function of 'func([]byte, int)' signature into
// CallableFunc type.
func FuncE(fn func() error) tengo.CallableFunc {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}

		return WrapError(fn()), nil
	}
}
