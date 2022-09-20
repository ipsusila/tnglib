package script

import (
	"context"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
)

func prepareExecutor(filename string, modules ...string) (string, Executor, error) {
	const id = "__default__"
	conf := DefaultConfig()
	conf.Modules = modules
	man := NewManager()
	exe := NewExecutor(man, Unlimited)
	if err := man.AddFile(id, filename, conf); err != nil {
		return id, nil, err
	}
	return id, exe, nil
}

// RunFile execute tengo script file using default context
func RunFile(filename string, modules ...string) error {
	id, exe, err := prepareExecutor(filename, modules...)
	if err != nil {
		return err
	}
	_, err = exe.Exec(id, nil)
	return err
}

// RunFileContext execute tengo script file with given context
func RunFileContext(ctx context.Context, filename string, modules ...string) error {
	id, exe, err := prepareExecutor(filename, modules...)
	if err != nil {
		return err
	}

	_, err = exe.ExecContext(ctx, id, nil)
	return err
}

// GetModuleMap from given modules name
func GetModuleMap(modules []string) *tengo.ModuleMap {
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
