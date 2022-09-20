package script

import (
	"path/filepath"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
)

// Config stores configuration
// for each script to be executed.
type Config struct {
	MaxExecutionTime string                 `json:"maxExecutionTime"`
	EnableFileImport bool                   `json:"enableFileImport"`
	ImportDir        string                 `json:"importDir"`
	MaxAllocs        int64                  `json:"maxAllocs"`
	MaxConstObjects  int                    `json:"maxConstObjects"`
	ImportFileExt    []string               `json:"importFileExt"`
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
	if err := tnglib.LoadConfigTo(filename, conf); err != nil {
		return nil, err
	}
	return conf, nil
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
func (c Config) ImportFileExtensions() []string {
	if len(c.ImportFileExt) == 0 {
		return []string{tengo.SourceFileExtDefault}
	}
	return c.ImportFileExt
}

// IsSourceFile return true if extension is registered as compiled extension
func (c Config) IsSourceFile(filename string) bool {
	ext := filepath.Ext(filename)
	extList := c.ImportFileExtensions()
	for _, v := range extList {
		if v == ext {
			return true
		}
	}
	return false
}
