package script

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
)

const (
	verCnt = 4
)

var (
	verVal = [verCnt]byte{'v', '2', '.', '0'} // as tengo version
)

type header struct {
	CompiledAt    time.Time
	Conf          Config
	Variables     map[string]interface{}
	GlobalIndexes map[string]int
}

type compiledBytecode struct {
	hdr      header
	bytecode *tengo.Bytecode
	globals  []tengo.Object
	filename string
}

// BytecodeFromSource compiles source file and then store it to bytecode
func BytecodeFromSource(filename string, conf *Config) (Entry, error) {
	bc := compiledBytecode{
		hdr: header{
			Conf: *conf,
		},
		filename: filename,
	}
	if err := bc.Recompile(); err != nil {
		return nil, err
	}
	return &bc, nil
}

// BytecodeFromFile loads compiled code from file
func BytecodeFromFile(filename string, conf *Config) (Entry, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	bc := compiledBytecode{
		hdr: header{
			Conf: *conf,
		},
	}
	if err := bc.Decode(fd); err != nil {
		return nil, err
	}

	return &bc, nil
}

func (b *compiledBytecode) Configuration() Config {
	return b.hdr.Conf
}
func (b *compiledBytecode) CompiledAt() time.Time {
	return b.hdr.CompiledAt
}
func (b *compiledBytecode) Age() time.Duration {
	return time.Since(b.hdr.CompiledAt)
}
func (b *compiledBytecode) Runnable() Runnable {
	clone := &bytecodeX{
		globalIndexes: b.hdr.GlobalIndexes,
		bytecode:      b.bytecode,
		globals:       make([]tengo.Object, len(b.globals)),
		maxAllocs:     b.hdr.Conf.MaxAllocs,
	}
	// copy global objects
	for idx, g := range b.globals {
		if g != nil {
			clone.globals[idx] = g
		}
	}
	return clone
}
func (b *compiledBytecode) Recompile() error {
	if b.filename == "" {
		// not supported
		return ErrBytecodeNotRecompilable
	}

	// recompfile source code
	input, err := os.ReadFile(b.filename)
	if err != nil {
		return err
	}
	if err := b.Compile(&b.hdr.Conf, input); err != nil {
		return err
	}
	return nil
}

func (b *compiledBytecode) AddMap(vars map[string]interface{}) {
	if len(b.hdr.Variables) == 0 {
		b.hdr.Variables = make(map[string]interface{})
	}
	for name, val := range vars {
		b.hdr.Variables[name] = val
	}
}
func (b *compiledBytecode) Add(name string, value interface{}) {
	if len(b.hdr.Variables) == 0 {
		b.hdr.Variables = make(map[string]interface{})
	}
	b.hdr.Variables[name] = value
}
func (b *compiledBytecode) SetMaxAllocs(n int64) {
	b.hdr.Conf.MaxAllocs = n
}
func (b *compiledBytecode) Remove(name string) bool {
	if _, ok := b.hdr.Variables[name]; !ok {
		return false
	}
	delete(b.hdr.Variables, name)
	return true
}

func (b *compiledBytecode) makeSymbolTable() (*tengo.SymbolTable, error) {
	var names []string
	for name := range b.hdr.Variables {
		names = append(names, name)
	}

	symbolTable := tengo.NewSymbolTable()
	builtinFuncs := tengo.GetAllBuiltinFunctions()
	for idx, fn := range builtinFuncs {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	globalsIndexes := make(map[string]int)
	for idx, name := range names {
		symbol := symbolTable.Define(name)
		if symbol.Index != idx {
			return nil, fmt.Errorf("wrong symbol index: %d != %d", idx, symbol.Index)
		}
		globalsIndexes[name] = idx
	}
	b.hdr.GlobalIndexes = globalsIndexes
	return symbolTable, nil
}

func (b *compiledBytecode) makeGlobals() ([]tengo.Object, error) {
	globals := make([]tengo.Object, tengo.GlobalsSize)
	maxIdx := -1
	for name, idx := range b.hdr.GlobalIndexes {
		obj, err := tengo.FromInterface(b.hdr.Variables[name])
		if err != nil {
			return nil, err
		}
		globals[idx] = obj
		if idx > maxIdx {
			maxIdx = idx
		}
	}
	if maxIdx >= 0 {
		globals = globals[:maxIdx+1]
	}
	return globals, nil
}
func (b *compiledBytecode) Compile(conf *Config, input []byte) error {
	// add variables from configuration
	b.AddMap(conf.InitVars)

	symbolTable, err := b.makeSymbolTable()
	if err != nil {
		return err
	}

	fileSet := parser.NewFileSet()
	srcFile := fileSet.AddFile("(main)", -1, len(input))
	p := parser.NewParser(srcFile, input, nil)
	file, err := p.ParseFile()
	if err != nil {
		return err
	}

	// get default import directory
	// 1. current
	// 2. working directory
	// 3. if filename is specified, same to the folder of the filename
	defImpDir := "."
	if dir, err := os.Getwd(); err == nil {
		defImpDir = dir
	}
	if b.filename != "" {
		defImpDir = filepath.Dir(b.filename)
	}

	modules := GetModuleMap(conf.Modules)
	c := tengo.NewCompiler(srcFile, symbolTable, nil, modules, nil)
	c.EnableFileImport(conf.EnableFileImport)
	c.SetImportDir(conf.ImportDirectory(filepath.Dir(defImpDir)))
	if err := c.Compile(file); err != nil {
		return err
	}

	// remove duplicates from constants
	bytecode := c.Bytecode()
	bytecode.RemoveDuplicates()

	// check the constant objects limit
	if conf.MaxConstObjects >= 0 {
		cnt := bytecode.CountObjects()
		if cnt > conf.MaxConstObjects {
			return fmt.Errorf("exceeding constant objects limit: %d", cnt)
		}
	}

	globals, err := b.makeGlobals()
	if err != nil {
		return err
	}

	b.hdr.Conf = *conf
	b.hdr.CompiledAt = time.Now()
	b.bytecode = bytecode
	b.globals = globals

	return nil
}

func (b *compiledBytecode) Encode(w io.Writer) error {
	// TODO: specify version
	return b.encode(w, verVal[:])
}
func (b *compiledBytecode) Decode(r io.Reader) error {
	// 1. read version
	ver := [verCnt]byte{}
	n, err := r.Read(ver[:])
	if err != nil {
		return err
	}
	if n != verCnt {
		return ErrInvalidBytecode
	}

	// TODO: verify version

	// 2. read header
	var hdr header
	dec := json.NewDecoder(r)
	if err := dec.Decode(&hdr); err != nil {
		return err
	}

	// 3. decode bytecode
	var bytecode tengo.Bytecode
	modules := GetModuleMap(b.hdr.Conf.Modules)
	if err := bytecode.Decode(r, modules); err != nil {
		return err
	}

	globals, err := b.makeGlobals()
	if err != nil {
		return err
	}

	b.hdr = hdr
	b.bytecode = &bytecode
	b.globals = globals
	b.filename = "" // not recomfileable source

	return nil
}

func (b *compiledBytecode) encode(w io.Writer, ver []byte) error {
	if b.bytecode == nil {
		return ErrBytecodeNotReady
	}
	if _, err := w.Write(ver); err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(b.hdr); err != nil {
		return err
	}
	return b.bytecode.Encode(w)
}
