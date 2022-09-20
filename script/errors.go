package script

import "errors"

// List of known errors
var (
	ErrScriptAlreadyRegistered = errors.New("script already registered")
	ErrScriptDoesNotExists     = errors.New("script does not exists")
	ErrUnsupportedConfigFormat = errors.New("unsupported configuration file format")
	ErrBytecodeNotReady        = errors.New("bytecode not ready")
	ErrInvalidBytecode         = errors.New("invalid bytecode")
	ErrBytecodeNotRecompilable = errors.New("can not recompile bytecode")
)
