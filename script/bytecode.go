package script

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
)

// file identifier
const (
	verCnt       = 4
	magicNo byte = 11
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
func BytecodeFromFile(filename string) (Entry, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	bc := compiledBytecode{}
	if err := bc.Decode(fd); err != nil {
		return nil, err
	}

	return &bc, nil
}

// BytecodeFromBytes convert byte array to byte code
func BytecodeFromBytes(input []byte) (Entry, error) {
	r := bytes.NewReader(input)
	bc := compiledBytecode{}
	if err := bc.Decode(r); err != nil {
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
	defImpDir, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	if dir, err := os.Getwd(); err == nil {
		defImpDir = filepath.Clean(dir)
	}
	if b.filename != "" {
		defImpDir = filepath.Dir(b.filename)
	}

	modules := GetModuleMap(conf.Modules)
	c := tengo.NewCompiler(srcFile, symbolTable, nil, modules, os.Stdout)
	c.SetImportFileExt(conf.ImportFileExtensions()...)
	c.EnableFileImport(conf.EnableFileImport)
	c.SetImportDir(conf.ImportDirectory(defImpDir))
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
	// 1. Get magic no
	mn := [1]byte{}
	n, err := r.Read(mn[:])
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrInvalidBytecode
	}

	// 2. Get version
	// TODO: verify version
	ver := [verCnt]byte{}
	n, err = r.Read(ver[:])
	if err != nil {
		return err
	}
	if n != verCnt {
		return ErrInvalidBytecode
	}

	// 3. Get number of bytes for Config stream
	cnt := make([]byte, 8)
	if _, err := r.Read(cnt); err != nil {
		return err
	}
	nb := binary.BigEndian.Uint64(cnt)
	if nb <= 0 {
		return ErrInvalidBytecode
	}
	js := make([]byte, nb)
	if _, err := r.Read(js); err != nil {
		return err
	}

	// 4. Get JSON config
	var hdr header
	if err := json.Unmarshal(js, &hdr); err != nil {
		return err
	}

	// 5. decode bytecode
	bytecode := &tengo.Bytecode{}
	modules := GetModuleMap(hdr.Conf.Modules)
	if err := bytecode.Decode(r, modules); err != nil {
		return err
	}

	// Create global variables
	globals, err := b.makeGlobals()
	if err != nil {
		return err
	}

	b.hdr = hdr
	b.bytecode = bytecode
	b.globals = globals
	b.filename = "" // not recompilable source

	return nil
}

func (b *compiledBytecode) encode(w io.Writer, ver []byte) error {
	if b.bytecode == nil {
		return ErrBytecodeNotReady
	}
	// 1. Write magic number
	if _, err := w.Write([]byte{magicNo}); err != nil {
		return err
	}

	// 2. Write version info
	if _, err := w.Write(ver); err != nil {
		return err
	}

	// 3. Write size of JSON stream and JSON config
	js, err := json.Marshal(b.hdr)
	if err != nil {
		return err
	}
	nb := uint64(len(js))
	cnt := make([]byte, 8)
	binary.BigEndian.PutUint64(cnt, nb)
	if _, err := w.Write(cnt); err != nil {
		return err
	}
	if _, err := w.Write(js); err != nil {
		return err
	}

	// 4. Write bytecode stream
	return b.bytecode.Encode(w)
}

// SaveTo saves bytecode to file
func (b *compiledBytecode) SaveTo(filename string) error {
	fw, err := os.Create(filename)
	if err != nil {
		return nil
	}
	defer fw.Close()

	return b.Encode(fw)
}
