package script

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/d5/tengo/v2"
	"golang.org/x/sync/semaphore"
)

const (
	Unlimited        = -1              // Unlimited number fo concurrency
	MaxExecutionTime = 5 * time.Minute // Default maximum execution time to 5 minutes
)

// Executor executes manages script
type Executor interface {
	Exec(id string, inpVars map[string]interface{}, outVars ...string) ([]*tengo.Variable, error)
	ExecContext(ctx context.Context, id string, inpVars map[string]interface{}, outVars ...string) ([]*tengo.Variable, error)
	NumInProgress() int
	WaitAll(ctx context.Context) error
}

type compiledExecutor struct {
	maxConcurrent int
	numInProgress atomic.Int32
	sem           *semaphore.Weighted
	man           Manager
}

// NewExecutor create script executor with limited concurrency.
// If set to <= 0, then there will be no limit.
func NewExecutor(man Manager, maxConcurrent int) Executor {
	e := compiledExecutor{
		man:           man,
		maxConcurrent: maxConcurrent,
	}
	if maxConcurrent > 0 {
		e.sem = semaphore.NewWeighted(int64(maxConcurrent))
	}
	return &e
}

func (c *compiledExecutor) NumInProgress() int {
	n := c.numInProgress.Load()
	return int(n)
}

func (c *compiledExecutor) WaitAll(ctx context.Context) error {
	if c.sem != nil {
		return c.sem.Acquire(ctx, int64(c.maxConcurrent))
	}
	return nil
}

func (c *compiledExecutor) Exec(id string, inpVars map[string]interface{}, outVars ...string) ([]*tengo.Variable, error) {
	entry := c.man.Entry(id)
	if entry == nil {
		return nil, fmt.Errorf("execute script id `%s`: %w", id, ErrScriptDoesNotExists)
	}

	conf := entry.Configuration()
	ctx, cancel := context.WithTimeout(context.TODO(), conf.MaxTimeout(MaxExecutionTime))
	defer cancel()

	return c.runContext(ctx, entry, inpVars, outVars...)
}

func (c *compiledExecutor) ExecContext(
	ctx context.Context,
	id string,
	inpVars map[string]interface{},
	outVars ...string,
) ([]*tengo.Variable, error) {

	entry := c.man.Entry(id)
	if entry == nil {
		return nil, fmt.Errorf("execute script id `%s`: %w", id, ErrScriptDoesNotExists)
	}
	// execute in context
	return c.runContext(ctx, entry, inpVars, outVars...)
}

func (c *compiledExecutor) runContext(
	ctx context.Context,
	entry Entry,
	inpVars map[string]interface{},
	outVars ...string,
) ([]*tengo.Variable, error) {

	// assign variables
	compiled := entry.Runnable()
	for name, val := range inpVars {
		if err := compiled.Set(name, val); err != nil {
			return nil, err
		}
	}

	var err error
	if c.maxConcurrent > 0 {
		err = c.runAsync(ctx, compiled)
	} else {
		c.numInProgress.Add(1)
		err = compiled.RunContext(ctx)
		c.numInProgress.Add(-1)
	}

	// check for error and return variable
	if err != nil {
		return nil, err
	}
	if len(outVars) > 0 {
		results := []*tengo.Variable{}
		for _, name := range outVars {
			results = append(results, compiled.Get(name))
		}
		return results, nil
	}
	return compiled.GetAll(), nil
}

func (c *compiledExecutor) runAsync(ctx context.Context, compiled Runnable) error {
	if c.sem == nil {
		return errors.New("bug: pool/limiter not defined")
	}
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	c.numInProgress.Add(1)
	errCh := make(chan error)
	go func() {
		defer func() {
			c.numInProgress.Add(-1)
			defer c.sem.Release(1)
		}()

		// execute the script
		if err := compiled.RunContext(ctx); err != nil {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	return <-errCh
}
