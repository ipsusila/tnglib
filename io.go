package tnglib

import (
	"errors"
	"io"
	"os"

	"github.com/d5/tengo/v2"
)

// list of writer
const (
	mWriter = "write"
	mReader = "read"
)

// standard io module
var ioModule = map[string]tengo.Object{
	"seek_set": &tengo.Int{Value: int64(io.SeekStart)},
	"seek_cur": &tengo.Int{Value: int64(io.SeekCurrent)},
	"seek_end": &tengo.Int{Value: int64(io.SeekEnd)},
	"stdout": &Writer{
		InterfaceImpl: InterfaceImpl{Name: "stdout"},
		Value:         os.Stdout,
	},
	"stderr": &Writer{
		InterfaceImpl: InterfaceImpl{Name: "stderr"},
		Value:         os.Stderr,
	},
	"stdin": &Reader{
		InterfaceImpl: InterfaceImpl{Name: "stdin"},
		Value:         os.Stdin,
	},
	"discard": &Writer{
		InterfaceImpl: InterfaceImpl{Name: "discard"},
		Value:         io.Discard,
	},
	"new_scanner": &tengo.UserFunction{
		Name:  "new_scanner",
		Value: newIoScanner,
	},
}

// IoFunc wrapper
type IoFunc struct {
	Fn *tengo.UserFunction
}

// NewIoFunc from tengo Object
func NewIoFunc(name string, obj tengo.Object) (*IoFunc, bool) {
	i := IoFunc{}
	ok := i.Set(name, obj)
	return &i, ok
}

func (i *IoFunc) validate(name string) error {
	if i.Fn == nil {
		return errors.New("tengo.UserFunction not assigned")
	}
	if i.Fn.Name != name {
		return errors.New("expected function name: " + name + ", got: " + i.Fn.Name)
	}
	return nil
}

func (i *IoFunc) callBIE(name string, data []byte) (int, error) {
	if err := i.validate(name); err != nil {
		return -1, err
	}
	arg := &tengo.Bytes{Value: data}
	ret, err := i.Fn.Call(arg)
	if err != nil {
		return -1, err
	}
	if err, ok := ret.(*tengo.Error); ok {
		return -1, errors.New(err.Value.String())
	}

	n, ok := ret.(*tengo.Int)
	if !ok {
		return 0, nil
	}
	return int(n.Value), nil
}

func (i *IoFunc) setMap(name string, m map[string]tengo.Object) bool {
	obj, ok := m[name]
	if !ok {
		return false
	}
	i.Fn, ok = obj.(*tengo.UserFunction)
	return ok
}

// Set user function with given name (if exist in the object)
func (i *IoFunc) Set(name string, o tengo.Object) bool {
	switch v := o.(type) {
	case *tengo.ImmutableMap:
		return i.setMap(name, v.Value)
	case *tengo.Map:
		return i.setMap(name, v.Value)
	case *tengo.UserFunction:
		if v.Name == name {
			i.Fn = v
			return true
		}
		return false
	default:
		return false
	}
}

// Write implement io.Writer
func (i *IoFunc) Write(data []byte) (int, error) {
	return i.callBIE(mWriter, data)
}

// Reader implement io.Reader
func (i *IoFunc) Read(data []byte) (int, error) {
	return i.callBIE(mReader, data)
}
