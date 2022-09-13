package tnglib

import (
	"fmt"
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
		return wrapError(fn(wr, da)), nil
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
			return wrapError(err), nil
		}
		return &tengo.String{Value: sb.String()}, nil
	}
}

// ArgToString convert tengo function call arguments to string.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgToString(args ...tengo.Object) (string, error) {
	if len(args) != 1 {
		return "", tengo.ErrWrongNumArguments
	}
	str, ok := tengo.ToString(args[0])
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	return str, nil
}

// ArgIToString convert tengo function call arguments to string.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToString(idx int, args ...tengo.Object) (string, error) {
	if idx >= len(args) {
		return "", tengo.ErrWrongNumArguments
	}
	str, ok := tengo.ToString(args[idx])
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "string(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return str, nil
}

// ArgsToStrings convert tengo function call arguments to string slice.
// If the argument count is less than minArg, it will return ErrWrongNumArguments
func ArgsToStrings(minArg int, args ...tengo.Object) ([]string, error) {
	if len(args) < minArg {
		return nil, tengo.ErrWrongNumArguments
	}
	items := []string{}
	for idx, arg := range args {
		filename, ok := tengo.ToString(arg)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("arg[%d]", idx),
				Expected: "string(compatible)",
				Found:    arg.TypeName(),
			}
		}
		items = append(items, filename)
	}
	return items, nil
}
