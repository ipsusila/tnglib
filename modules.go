package tnglib

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/registry"
)

// module registry
var modsLib = registry.NewImmutableSyncMapRegistry[string, map[string]tengo.Object]()

func init() {
	RegisterModule("text/template", textTplModule)
	RegisterModule("html/template", htmlTplModule)
	RegisterModule("io", ioModule)
	RegisterModule("context", ctxModule)
}

// RegisterModule with given name
func RegisterModule(name string, mod map[string]tengo.Object) {
	if err := modsLib.Register(name, mod); err != nil {
		panic(err)
	}
}

// GetModuleMap returns the module map that includes all modules
// for the given module names. If module does not exist, it will panic.
func GetModuleMap(names ...string) *tengo.ModuleMap {
	modules := tengo.NewModuleMap()
	for _, name := range names {
		if mod, err := modsLib.Entry(name); err == nil {
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

// AllModuleNames return all registered modules
func AllModuleNames() []string {
	return modsLib.Keys()
}
