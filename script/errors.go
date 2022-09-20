package script

import "errors"

// List of known errors
var (
	ErrEntryAlreadyRegistered  = errors.New("entry already registered")
	ErrEntryDoesNotExists      = errors.New("entry does not exists")
	ErrUnsupportedConfigFormat = errors.New("unsupported configuration file format")
	ErrBytecodeNotReady        = errors.New("bytecode not ready")
	ErrInvalidBytecode         = errors.New("invalid bytecode")
	ErrBytecodeNotRecompilable = errors.New("can not recompile bytecode")
)
