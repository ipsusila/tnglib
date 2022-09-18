package tnglib_test

import (
	"os"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
)

func runTengoScript(filename string, modules ...string) error {
	script, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// create script, load all std modules and tnglib modules
	s := tengo.NewScript([]byte(script))
	mod := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	lib := tnglib.GetModuleMap(modules...)
	mod.AddMap(lib)
	s.SetImports(mod)

	// compile the source
	c, err := s.Compile()
	if err != nil {
		return err
	}

	// run the compiled bytecode
	// a compiled bytecode 'c' can be executed multiple times without re-compiling it
	return c.Run()
}
