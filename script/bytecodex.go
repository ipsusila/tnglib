package script

import (
	"context"
	"fmt"
	"sync"

	"github.com/d5/tengo/v2"
)

// Taken from https://github.com/d5/tengo/blob/v2.12.2/script.go#L13

// bytecodeX is a compiled instance of the user script. Use Script.Compile() to
// create bytecodeX object.
type bytecodeX struct {
	globalIndexes map[string]int // global symbol name to index
	bytecode      *tengo.Bytecode
	globals       []tengo.Object
	maxAllocs     int64
	lock          sync.RWMutex
}

// ensure implement runnable
var _ Runnable = (*bytecodeX)(nil)

// Run executes the compiled script in the virtual machine.
func (c *bytecodeX) Run() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	v := tengo.NewVM(c.bytecode, c.globals, c.maxAllocs)
	return v.Run()
}

// RunContext is like Run but includes a context.
func (c *bytecodeX) RunContext(ctx context.Context) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	v := tengo.NewVM(c.bytecode, c.globals, c.maxAllocs)
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case string:
					ch <- fmt.Errorf(e)
				case error:
					ch <- e
				default:
					ch <- fmt.Errorf("unknown panic: %v", e)
				}
			}
		}()
		ch <- v.Run()
	}()

	select {
	case <-ctx.Done():
		v.Abort()
		<-ch
		err = ctx.Err()
	case err = <-ch:
	}
	return
}

// IsDefined returns true if the variable name is defined (has value) before or
// after the execution.
func (c *bytecodeX) IsDefined(name string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, ok := c.globalIndexes[name]
	if !ok {
		return false
	}
	v := c.globals[idx]
	if v == nil {
		return false
	}
	return v != tengo.UndefinedValue
}

// Get returns a variable identified by the name.
func (c *bytecodeX) Get(name string) *tengo.Variable {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value := tengo.UndefinedValue
	if idx, ok := c.globalIndexes[name]; ok {
		value = c.globals[idx]
		if value == nil {
			value = tengo.UndefinedValue
		}
	}
	return c.newVariable(name, value)
}
func (c *bytecodeX) newVariable(name string, value tengo.Object) *tengo.Variable {
	v, _ := tengo.NewVariable(name, value)
	return v
}

// GetAll returns all the variables that are defined by the compiled script.
func (c *bytecodeX) GetAll() []*tengo.Variable {
	c.lock.RLock()
	defer c.lock.RUnlock()

	var vars []*tengo.Variable
	for name, idx := range c.globalIndexes {
		value := c.globals[idx]
		if value == nil {
			value = tengo.UndefinedValue
		}
		vars = append(vars, c.newVariable(name, value))
	}
	return vars
}

// Set replaces the value of a global variable identified by the name. An error
// will be returned if the name was not defined during compilation.
func (c *bytecodeX) Set(name string, value interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	obj, err := tengo.FromInterface(value)
	if err != nil {
		return err
	}
	idx, ok := c.globalIndexes[name]
	if !ok {
		return fmt.Errorf("'%s' is not defined", name)
	}
	c.globals[idx] = obj
	return nil
}
