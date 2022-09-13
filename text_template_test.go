package tnglib_test

import (
	"os"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/stretchr/testify/assert"
)

func runTengoScript(filename string) error {
	script, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// create script, load all std modules and tnglib modules
	s := tengo.NewScript([]byte(script))
	mod := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	lib := tnglib.GetModuleMap("io", "text/template")
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
func TestTextTemplate(t *testing.T) {
	err := runTengoScript("_testdata/template1.tengo")
	assert.NoError(t, err)
}

func BenchmarkTextTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := runTengoScript("_testdata/template1.tengo")
		if err != nil {
			b.Log(err)
		}
	}
}
