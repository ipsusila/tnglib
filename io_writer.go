package tnglib

import (
	"io"

	"github.com/d5/tengo/v2"
)

// Writer represents a user function.
type Writer struct {
	InterfaceImpl
	Value io.Writer
}

// Copy returns a copy of the type.
func (o *Writer) Copy() tengo.Object {
	return &Writer{
		InterfaceImpl: InterfaceImpl{
			Name: o.Name,
		},
		Value: o.Value,
	}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Writer) Equals(v tengo.Object) bool {
	an, ok := v.(*Writer)
	if !ok {
		return false
	}
	return o.Name == an.Name && o.Value == an.Value
}

// Call invokes a user function.
func (o *Writer) Call(args ...tengo.Object) (tengo.Object, error) {
	data, err := ArgToByteSlice(args...)
	if err != nil {
		return nil, err
	}

	n, err := o.Value.Write(data)
	if err != nil {
		return WrapError(err), nil
	}

	return &tengo.Int{Value: int64(n)}, nil
}

// CanCall returns whether the Object can be Called.
func (o *Writer) IsFalsy() bool {
	return o.Value == nil || o.Name == ""
}
