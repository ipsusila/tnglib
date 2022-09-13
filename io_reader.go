package tnglib

import (
	"io"

	"github.com/d5/tengo/v2"
)

// Reader represents a user function.
type Reader struct {
	InterfaceImpl
	Value io.Reader
}

// Copy returns a copy of the type.
func (o *Reader) Copy() tengo.Object {
	return &Reader{
		InterfaceImpl: InterfaceImpl{
			Name: o.Name,
		},
		Value: o.Value,
	}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Reader) Equals(v tengo.Object) bool {
	if v == nil || o == v {
		return o == v
	}
	an, ok := v.(*Reader)
	if !ok {
		return false
	}
	return o.Name == an.Name && o.Value == an.Value
}

// Call invokes a user function.
func (o *Reader) Call(args ...tengo.Object) (tengo.Object, error) {
	data, err := ArgToByteSlice(args...)
	if err != nil {
		return nil, err
	}

	n, err := o.Value.Read(data)
	if err != nil {
		return WrapError(err), nil
	}

	return &tengo.Int{Value: int64(n)}, nil
}

// CanCall returns whether the Object can be Called.
func (o *Reader) IsFalsy() bool {
	return o.Value == nil || o.Name == ""
}
