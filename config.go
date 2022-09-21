package tnglib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjson/hjson-go/v4"
)

// ErrUnsupportedConfigFormat return unsupported config file format
var ErrUnsupportedConfigFormat = errors.New("unsupported configuration file format")

// Ptr for pointer constraints
type Ptr[T any] interface {
	*T
}

// LoadConfigTo load configuration from given file.
// Supported formats: json, hjson
func LoadConfigTo[T any, P Ptr[T]](filename string, out *T) error {
	ext := strings.ToLower(filepath.Ext(filename))
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	switch ext {
	case ".json":
		dec := json.NewDecoder(fd)
		return dec.Decode(out)
	case ".hjson":
		var tout map[string]interface{}
		data, err := io.ReadAll(fd)
		if err != nil {
			return err
		}
		if err := hjson.Unmarshal(data, &tout); err != nil {
			return err
		}

		// marshal to JSON
		js, err := json.Marshal(tout)
		if err != nil {
			return err
		}
		return json.Unmarshal(js, out)
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
