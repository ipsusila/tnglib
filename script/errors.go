package script

import "errors"

// List of known errors
var (
	ErrScriptAlreadyRegistered = errors.New("script already registered")
	ErrScriptDoesNotExists     = errors.New("script does not exists")
	ErrUnsupportedConfigFormat = errors.New("unsupported configuration file format")
)
