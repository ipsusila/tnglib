package script

import (
	"context"

	"github.com/d5/tengo/v2"
)

// Runnable code
type Runnable interface {
	Get(name string) *tengo.Variable
	GetAll() []*tengo.Variable
	IsDefined(name string) bool
	Run() error
	RunContext(ctx context.Context) (err error)
	Set(name string, value interface{}) error
}
