package tnglib

import (
	"github.com/d5/tengo/v2"
)

// WrapError convert error to tengo object or true if err is nil.
func WrapError(err error) tengo.Object {
	if err == nil {
		return tengo.TrueValue
	}
	return &tengo.Error{Value: &tengo.String{Value: err.Error()}}
}
