package script

import (
	"os"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
)

// Entry properties
type Entry interface {
	ID() string
	Filename() string
	Configuration() Config
	CompiledAt() time.Time
	Age() time.Duration
	Compiled() *tengo.Compiled
}

type scriptEntry struct {
	id          string
	srcFilename string
	conf        *Config
	compiledAt  time.Time
	compiled    *tengo.Compiled
}

func (e *scriptEntry) ID() string {
	return e.id
}
func (e *scriptEntry) Filename() string {
	return e.srcFilename
}
func (e *scriptEntry) Configuration() Config {
	if e.conf == nil {
		return Config{}
	}
	return *e.conf
}
func (e *scriptEntry) CompiledAt() time.Time {
	return e.compiledAt
}
func (e *scriptEntry) Age() time.Duration {
	return time.Since(e.compiledAt)
}
func (e *scriptEntry) Compiled() *tengo.Compiled {
	if e.compiled == nil {
		return nil
	}
	return e.compiled.Clone()
}

func (e *scriptEntry) loadAndCompile() error {
	// read source code
	srcContent, err := os.ReadFile(e.srcFilename)
	if err != nil {
		return err
	}

	// load modules
	var mod *tengo.ModuleMap
	if len(e.conf.Modules) > 0 {
		mod = stdlib.GetModuleMap(e.conf.Modules...)
		mod.AddMap(tnglib.GetModuleMap(e.conf.Modules...))
	} else {
		mod = stdlib.GetModuleMap(stdlib.AllModuleNames()...)
		mod.AddMap(tnglib.GetModuleMap(tnglib.AllModuleNames()...))
	}

	// create script and import modules
	sc := tengo.NewScript([]byte(srcContent))
	sc.SetImports(mod)
	sc.EnableFileImport(e.conf.EnableFileImport)
	sc.SetImportDir(e.conf.ImportDirectory(e.srcFilename))
	if e.conf.MaxAllocs > 0 {
		sc.SetMaxAllocs(e.conf.MaxAllocs)
	}
	if e.conf.MaxConstObjects > 0 {
		sc.SetMaxConstObjects(e.conf.MaxConstObjects)
	}

	// add variabels
	for name, val := range e.conf.InitVars {
		if err := sc.Add(name, val); err != nil {
			return err
		}
	}

	// compile
	compiled, err := sc.Compile()
	if err != nil {
		return err
	}

	// store script and compiled
	e.compiled = compiled
	e.compiledAt = time.Now()

	return nil
}
