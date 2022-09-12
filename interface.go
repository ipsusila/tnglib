package tnglib

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

// InterfaceImpl for wraping interface
type InterfaceImpl struct {
	Name string
}

// TypeName returns the name of the type.
func (o *InterfaceImpl) TypeName() string {
	return "interface:" + o.Name
}

func (o *InterfaceImpl) String() string {
	return "<interface>"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *InterfaceImpl) BinaryOp(_ token.Token, _ tengo.Object) (tengo.Object, error) {
	return nil, tengo.ErrInvalidOperator
}

// IndexGet returns an element at a given index.
func (o *InterfaceImpl) IndexGet(_ tengo.Object) (res tengo.Object, err error) {
	return nil, tengo.ErrNotIndexable
}

// IndexSet sets an element at a given index.
func (o *InterfaceImpl) IndexSet(_, _ tengo.Object) (err error) {
	return tengo.ErrNotIndexAssignable
}

// Iterate returns an iterator.
func (o *InterfaceImpl) Iterate() tengo.Iterator {
	return nil
}

// CanIterate returns whether the tengo.Object can be Iterated.
func (o *InterfaceImpl) CanIterate() bool {
	return false
}

// CanCall returns whether the tengo.Object can be Called.
func (o *InterfaceImpl) CanCall() bool {
	return true
}

/*
// Call takes an arbitrary number of arguments and returns a return value
// and/or an error.
func (o *InterfaceImpl) Call(_ ...tengo.Object) (ret tengo.Object, err error) {
	return nil, nil
}

// Copy returns a copy of the type.
func (o *InterfaceImpl) Copy() tengo.Object {
	return nil
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *InterfaceImpl) Equals(_ tengo.Object) bool {
	return false
}

// IsFalsy returns true if the value of the type is falsy.
func (o *InterfaceImpl) IsFalsy() bool {
	return false
}
*/
