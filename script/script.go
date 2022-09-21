package script

import (
	"context"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
)

func prepareRunnable(filename string, modules ...string) (Runnable, error) {
	conf := DefaultConfig()
	conf.Modules = append(conf.Modules, modules...)
	e, err := BytecodeFromFile(filename, conf)
	if err != nil {
		return nil, err
	}
	return e.Runnable(), nil
}

// RunFile execute tengo script file using default context
func RunFile(filename string, modules ...string) error {
	r, err := prepareRunnable(filename, modules...)
	if err != nil {
		return err
	}
	return r.Run()
}

// RunFileContext execute tengo script file with given context
func RunFileContext(ctx context.Context, filename string, modules ...string) error {
	r, err := prepareRunnable(filename, modules...)
	if err != nil {
		return err
	}
	return r.RunContext(ctx)
}

// GetImportableModuleMap from given modules name.
// Search module from stdlib and tnglib
func GetImportableModuleMap(modules []string) *tengo.ModuleMap {
	var mod *tengo.ModuleMap
	if len(modules) > 0 {
		mod = stdlib.GetModuleMap(modules...)
		mod.AddMap(tnglib.GetModuleMap(modules...))
	} else {
		mod = stdlib.GetModuleMap(stdlib.AllModuleNames()...)
		mod.AddMap(tnglib.GetModuleMap(tnglib.AllModuleNames()...))
	}
	return mod
}
