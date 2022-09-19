package tnglib

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrUnsupportedConfigFormat return unsupported config file format
var ErrUnsupportedConfigFormat = errors.New("unsupported configuration file format")

// Ptr for pointer constraints
type Ptr[T any] interface {
	*T
}

// LoadConfigTo load configuration from given file.
// Supported formats: json
func LoadConfigTo[T any, P Ptr[T]](filename string, out *T) error {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		fd, err := os.Open(filename)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(fd)
		return dec.Decode(out)
	}
	return fmt.Errorf("config file `%s`: %w", filename, ErrUnsupportedConfigFormat)
}

// LoadConfig load configuration from given file.
// Supported formats: json
func LoadConfig[T any, P Ptr[T]](filename string) (*T, error) {
	var out T
	err := LoadConfigTo(filename, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
