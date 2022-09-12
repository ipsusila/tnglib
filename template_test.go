package tnglib_test

import (
	"os"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	"github.com/stretchr/testify/assert"
)

func runTengoScript(t *testing.T, filename string) {
	script, err := os.ReadFile(filename)
	assert.NoError(t, err)

	// create script
	s := tengo.NewScript([]byte(script))
	mod := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	lib := tnglib.GetModuleMap("io", "text/template")
	mod.AddMap(lib)
	s.SetImports(mod)

	// compile the source
	c, err := s.Compile()
	assert.NoError(t, err)

	// run the compiled bytecode
	// a compiled bytecode 'c' can be executed multiple times without re-compiling it
	if err := c.Run(); err != nil {
		assert.Error(t, err)
		panic(err)
	}
}
func TestTextTemplate(t *testing.T) {
	runTengoScript(t, "_testdata/template1.tengo")
}
