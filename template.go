package tnglib

import (
	"io"

	"github.com/d5/tengo/v2"
)

func tplExecute(fn func(w io.Writer, data any) error) tengo.CallableFunc {
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
