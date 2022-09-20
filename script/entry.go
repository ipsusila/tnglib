package script

import (
	"os"
	"time"

	"github.com/d5/tengo/v2"
)

// Entry properties
type Entry interface {
	Configuration() Config
	CompiledAt() time.Time
	Age() time.Duration
	Runnable() Runnable
	Recompile() error
}

type scriptEntry struct {
	srcFilename string
	conf        *Config
	compiledAt  time.Time
	compiled    *tengo.Compiled
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

func (e *scriptEntry) Runnable() Runnable {
	if e.compiled == nil {
		return nil
	}
	return e.compiled.Clone()
}
func (e *scriptEntry) Recompile() error {
	return e.loadAndCompileSrcript()
}
func (e *scriptEntry) loadAndCompileSrcript() error {
	// read source code
	srcContent, err := os.ReadFile(e.srcFilename)
	if err != nil {
		return err
	}

	// create script and import modules
	sc := tengo.NewScript([]byte(srcContent))
	sc.SetImports(GetModuleMap(e.conf.Modules))
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

func loadAndCompileSrcript(srcFilename string, conf *Config) (Entry, error) {
	// create new entry
	e := scriptEntry{
		srcFilename: srcFilename,
		conf:        conf,
	}
	if err := e.loadAndCompileSrcript(); err != nil {
		return nil, err
	}

	return &e, nil
}
