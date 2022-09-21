package tnglib

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/d5/tengo/v2"
)

// FuncAWARE convert any function(io.Writer, interface{}) error to tengo.CallableFunc
func FuncAWARE(fn func(w io.Writer, data any) error) tengo.CallableFunc {
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

// FuncAWAREs convert any function(io.Writer, interface{}) error to tengo.CallableFunc
func FuncAWAREs(fn func(w io.Writer, data any) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}

		var sb strings.Builder
		da := tengo.ToInterface(args[0])
		if err := fn(&sb, da); err != nil {
			return WrapError(err), nil
		}
		s := sb.String()
		if len(s) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		return &tengo.String{Value: s}, nil
	}
}

// FuncARYs transform a function of 'func() []byte' signature into
// CallableFunc type.
func FuncARYs(fn func() []byte) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		b := fn()
		if len(b) > tengo.MaxBytesLen {
			return nil, tengo.ErrBytesLimit
		}
		return &tengo.Bytes{Value: b}, nil
	}
}

// FuncAYIR transform a function of 'func([]byte, int)' signature into
// CallableFunc type.
func FuncAYIR(fn func([]byte, int)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
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

		return tengo.UndefinedValue, nil
	}
}

// FuncATBR transform a function of 'func() (time.Time, bool)' signature into
// CallableFunc type.
func FuncATBR(fn func() (time.Time, bool)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		tm, ok := fn()
		return &tengo.ImmutableMap{
			Value: map[string]tengo.Object{
				"time": &tengo.Time{Value: tm},
				"ok":   BoolObject(ok),
			},
		}, nil
	}
}

// FuncAIRI transform a function of 'func(any) any' signature into
// CallableFunc type.
func FuncAIRI(fn func(any) any) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		arg := tengo.ToInterface(args[0])
		res := fn(arg)

		return tengo.FromInterface(res)
	}
}

// FuncAISRS transform a function of 'func(int, string) string' signature into
// CallableFunc type.
func FuncAISRS(fn func(int, string) string) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		i, err := ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		s, err := ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}
		res := fn(i, s)
		if len(res) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		return &tengo.String{Value: res}, nil
	}
}

func FuncAISARSAE(fn func(int, string, interface{}) (string, []interface{}, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}
		i, err := ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		s, err := ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}
		v := tengo.ToInterface(args[2])

		rs, ra, err := fn(i, s, v)
		if err != nil {
			return WrapError(err), nil
		}
		if len(rs) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		r0 := &tengo.String{Value: rs}
		r1, err := tengo.FromInterface(ra)
		if err != nil {
			return WrapError(err), nil
		}
		return &tengo.Array{Value: []tengo.Object{r0, r1}}, nil
	}
}

func FuncASVRSAE(fn func(string, ...interface{}) (string, []interface{}, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		iargs := []interface{}{}
		for i := 1; i < len(args); i++ {
			iargs = append(iargs, tengo.ToInterface(args[i]))
		}
		rs, ra, err := fn(s, iargs...)
		if err != nil {
			return WrapError(err), nil
		}
		if len(rs) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		r0 := &tengo.String{Value: rs}
		r1, err := tengo.FromInterface(ra)
		if err != nil {
			return WrapError(err), nil
		}
		return &tengo.Array{Value: []tengo.Object{r0, r1}}, nil
	}
}
func FuncASARSAE(fn func(string, interface{}) (string, []interface{}, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		s, err := ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		v := tengo.ToInterface(args[1])

		rs, ra, err := fn(s, v)
		if err != nil {
			return WrapError(err), nil
		}
		if len(rs) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		r0 := &tengo.String{Value: rs}
		r1, err := tengo.FromInterface(ra)
		if err != nil {
			return WrapError(err), nil
		}
		return &tengo.Array{Value: []tengo.Object{r0, r1}}, nil
	}
}

// FuncACRE transform a function of 'func(ctx) error' signature into
// CallableFunc type.
func FuncACRE(fn func(context.Context) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		ctx, err := ArgToContext(args...)
		if err != nil {
			return nil, err
		}
		err = fn(ctx.Ctx)
		return WrapError(err), nil
	}
}

// FuncASRI transform a function of 'func(string) int' signature
// into CallableFunc type.
func FuncASRI(fn func(string) int) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		return &tengo.Int{Value: int64(fn(s))}, nil
	}
}

// FuncARSsE transform a function of 'func() ([]string, error)' signature
// into CallableFunc type.
func FuncARSsE(fn func() ([]string, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		res, err := fn()
		if err != nil {
			return WrapError(err), nil
		}
		arr := &tengo.Array{}
		for _, r := range res {
			if len(r) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			arr.Value = append(arr.Value, &tengo.String{Value: r})
		}
		return arr, nil
	}
}

func FuncABRE(fn func(bool) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		v, err := ArgToBool(args...)
		if err != nil {
			return nil, err
		}
		return WrapError(fn(v)), nil
	}
}

func FuncADRE(fn func(time.Duration) error) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		v, cerr, err := ArgToDuration(args...)
		if err != nil {
			return nil, err
		} else if cerr != nil {
			return WrapError(cerr), nil
		}
		return WrapError(fn(v)), nil
	}
}
