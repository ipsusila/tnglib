package tnglib

import (
	"github.com/d5/tengo/v2"
)

// BuiltinModules are builtin type standard library modules.
var BuiltinModules = map[string]map[string]tengo.Object{
	"text/template": txtTplModule,
	"io":            ioModule,
}

// GetModuleMap returns the module map that includes all modules
// for the given module names.
func GetModuleMap(names ...string) *tengo.ModuleMap {
	modules := tengo.NewModuleMap()
	for _, name := range names {
		if mod := BuiltinModules[name]; mod != nil {
			modules.AddBuiltinModule(name, mod)
		}
		/*
			if mod := SourceModules[name]; mod != "" {
				modules.AddSourceModule(name, []byte(mod))
			}
		*/
	}
	return modules
}
