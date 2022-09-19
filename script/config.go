package script

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config stores configuration
// for each script to be executed.
type Config struct {
	MaxExecutionTime string                 `json:"maxExecutionTime"`
	EnableFileImport bool                   `json:"enableFileImport"`
	ImportDir        string                 `json:"importDir"`
	MaxAllocs        int64                  `json:"maxAllocs"`
	MaxConstObjects  int                    `json:"maxConstObjects"`
	Modules          []string               `json:"modules"` // empty means all modules
	InitVars         map[string]interface{} `json:"initVars"`
}

// DefaultConfig create default configuration for script
func DefaultConfig() *Config {
	return &Config{
		MaxAllocs:        -1,
		MaxConstObjects:  -1,
		EnableFileImport: true,
	}
}

// NewFileConfig load configuration from given file.
// Supported formats: json
func NewFileConfig(filename string) (*Config, error) {
	conf := DefaultConfig()
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		fd, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		dec := json.NewDecoder(fd)
		if err := dec.Decode(conf); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("config file `%s`: %w", filename, ErrUnsupportedConfigFormat)
}

// MaxTimeout return maximum execution time in time.Duration
func (c Config) MaxTimeout(def time.Duration) time.Duration {
	dur, err := time.ParseDuration(c.MaxExecutionTime)
	if err != nil {
		return def
	}
	return dur
}

// ImportDirectory for given srcFilename
func (c Config) ImportDirectory(srcFilename string) string {
	if c.ImportDir != "" {
		return c.ImportDir
	}
	return filepath.Dir(srcFilename)
}
