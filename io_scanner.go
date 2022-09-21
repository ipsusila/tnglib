package tnglib

import (
	"bufio"
	"io"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func makeIoScanner(s *bufio.Scanner) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// text() => string
			"text": &tengo.UserFunction{
				Name:  "text",
				Value: stdlib.FuncARS(s.Text),
			},
			// scan() => bool
			"scan": &tengo.UserFunction{
				Name:  "scan",
				Value: stdlib.FuncARB(s.Scan),
			},
			// bytes() => []byte
			"bytes": &tengo.UserFunction{
				Name:  "bytes",
				Value: FuncARYs(s.Bytes),
			},
			// buffer([]byte, int)
			"buffer": &tengo.UserFunction{
				Name:  "buffer",
				Value: FuncAYIR(s.Buffer),
			},
			// err() => error
			"err": &tengo.UserFunction{
				Name:  "err",
				Value: stdlib.FuncARE(s.Err),
			},
		},
	}
}

func newIoScanner(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	// get reader interface
	var rd io.Reader
	if v, ok := args[0].(*Reader); ok {
		rd = v.Value
	} else {
		rd, ok = NewIoFunc(mReader, args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "io.Reader",
				Found:    args[0].TypeName(),
			}
		}
	}
	s := bufio.NewScanner(rd)
	return makeIoScanner(s), nil
}
