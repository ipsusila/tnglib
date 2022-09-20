package script

import (
	"time"
)

// Entry properties
type Entry interface {
	Configuration() Config
	CompiledAt() time.Time
	Age() time.Duration
	Runnable() Runnable
	Recompile() error
	SaveTo(filename string) error
}
