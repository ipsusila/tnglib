package script

import (
	"fmt"
	"sync"
	"time"
)

// Manager manages compiled script
type Manager interface {
	Exists(id string) bool
	AddEntry(id string, e Entry) error
	Add(id string, data []byte) error
	AddFile(id, filename string, conf *Config) error
	Recompile(id string) error
	Entry(id string) Entry
	Remove(id string) bool
	Clear()
	Age(id string) (time.Duration, bool)
}

type manager struct {
	entries sync.Map
}

// NewManager creates compiled script manager
func NewManager() Manager {
	return &manager{}
}

func (m *manager) Exists(id string) bool {
	_, found := m.entries.Load(id)
	return found
}

func (m *manager) AddEntry(id string, e Entry) error {
	if m.Exists(id) {
		return ErrEntryAlreadyRegistered
	}
	m.entries.Store(id, e)
	return nil
}

func (m *manager) Add(id string, data []byte) error {
	if m.Exists(id) {
		return ErrEntryAlreadyRegistered
	}
	e, err := BytecodeFromBytes(data)
	if err != nil {
		return err
	}
	// store entries
	m.entries.Store(id, e)

	return nil
}

// Add new script to manager
func (m *manager) AddFile(id, filename string, conf *Config) error {
	// if configuration not specified,
	// use default configuration
	if conf == nil {
		conf = DefaultConfig()
	}
	if m.Exists(id) {
		return ErrEntryAlreadyRegistered
	}

	var e Entry
	var err error
	isSource := conf.IsSourceFile(filename)
	if isSource {
		e, err = BytecodeFromSource(filename, conf)
	} else {
		e, err = BytecodeFromFile(filename)
	}
	if err != nil {
		return err
	}

	// store entries
	m.entries.Store(id, e)
	return nil
}

// Recompile registered script
func (m *manager) Recompile(id string) error {
	v, found := m.entries.Load(id)
	if !found {
		return fmt.Errorf("recompile script with id `%s`: %w", id, ErrEntryDoesNotExists)
	}
	e := v.(Entry)
	return e.Recompile()
}

func (m *manager) Entry(id string) Entry {
	v, found := m.entries.Load(id)
	if !found {
		return nil
	}
	return v.(Entry)
}

func (m *manager) Remove(id string) bool {
	exists := m.Exists(id)
	m.entries.Delete(id)
	return exists
}

func (m *manager) Clear() {
	keys := []any{}
	m.entries.Range(func(key, value any) bool {
		keys = append(keys, key)
		return true
	})
	for _, key := range keys {
		m.entries.Delete(key)
	}
}
func (m *manager) Age(id string) (time.Duration, bool) {
	v, found := m.entries.Load(id)
	if !found {
		return 0, false
	}
	return v.(Entry).Age(), true
}
