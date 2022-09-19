package script

import (
	"fmt"
	"sync"
	"time"
)

// Manager manages compiled script
type Manager interface {
	Add(id, srcFilename string, conf *Config) error
	Recompile(id string) error
	Entry(id string) Entry
	Remove(id string)
	Clear()
	Age(id string) (time.Duration, bool)
}

type compiledManager struct {
	entries sync.Map
}

// NewManager creates compiled script manager
func NewManager() Manager {
	return &compiledManager{}
}

// Add new script to manager
func (c *compiledManager) Add(id, srcFilename string, conf *Config) error {
	// if configuration not specified,
	// use default configuration
	if conf == nil {
		conf = DefaultConfig()
	}
	_, found := c.entries.Load(id)
	if found {
		return fmt.Errorf("adding script with id `%s`: %w", id, ErrScriptAlreadyRegistered)
	}

	// create new script entry
	e := scriptEntry{
		id:          id,
		srcFilename: srcFilename,
		conf:        conf,
	}
	if err := e.loadAndCompile(); err != nil {
		return err
	}

	// store entries
	c.entries.Store(id, &e)
	return nil
}

// Recompile registered script
func (c *compiledManager) Recompile(id string) error {
	v, found := c.entries.Load(id)
	if !found {
		return fmt.Errorf("recompile script with id `%s`: %w", id, ErrScriptDoesNotExists)
	}
	e := v.(*scriptEntry)
	return e.loadAndCompile()
}

func (c *compiledManager) Entry(id string) Entry {
	v, found := c.entries.Load(id)
	if !found {
		return nil
	}
	return v.(*scriptEntry)
}

func (c *compiledManager) Remove(id string) {
	c.entries.Delete(id)
}

func (c *compiledManager) Clear() {
	keys := []any{}
	c.entries.Range(func(key, value any) bool {
		keys = append(keys, key)
		return true
	})
	for _, key := range keys {
		c.entries.Delete(key)
	}
}
func (c *compiledManager) Age(id string) (time.Duration, bool) {
	v, found := c.entries.Load(id)
	if !found {
		return 0, false
	}
	return time.Since(v.(*scriptEntry).compiledAt), true
}
