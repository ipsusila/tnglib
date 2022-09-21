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

type executor struct {
	maxConcurrent int
	numInProgress atomic.Int32
	sem           *semaphore.Weighted
	man           Manager
}

// NewExecutor create script executor with limited concurrency.
// If set to <= 0, then there will be no limit.
func NewExecutor(man Manager, maxConcurrent int) Executor {
	e := executor{
		man:           man,
		maxConcurrent: maxConcurrent,
	}
	if maxConcurrent > 0 {
		e.sem = semaphore.NewWeighted(int64(maxConcurrent))
	}
	return &e
}

func (e *executor) NumInProgress() int {
	n := e.numInProgress.Load()
	return int(n)
}

func (e *executor) WaitAll(ctx context.Context) error {
	if e.sem != nil {
		return e.sem.Acquire(ctx, int64(e.maxConcurrent))
	}
	return nil
}

func (e *executor) Exec(id string,
	inpVars map[string]interface{},
	outVars ...string,
) ([]*tengo.Variable, error) {
	entry := e.man.Entry(id)
	if entry == nil {
		return nil, fmt.Errorf("execute script id `%s`: %w", id, ErrEntryDoesNotExists)
	}

	conf := entry.Configuration()
	ctx, cancel := context.WithTimeout(context.TODO(), conf.MaxExecutionTime.Duration)
	defer cancel()

	return e.runContext(ctx, entry, inpVars, outVars...)
}

func (e *executor) ExecContext(
	ctx context.Context,
	id string,
	inpVars map[string]interface{},
	outVars ...string,
) ([]*tengo.Variable, error) {

	entry := e.man.Entry(id)
	if entry == nil {
		return nil, fmt.Errorf("execute script id `%s`: %w", id, ErrEntryDoesNotExists)
	}
	// execute in context
	return e.runContext(ctx, entry, inpVars, outVars...)
}

func (e *executor) runContext(
	ctx context.Context,
	entry Entry,
	inpVars map[string]interface{},
	outVars ...string,
) ([]*tengo.Variable, error) {

	// assign variables
	r := entry.Runnable()
	for name, val := range inpVars {
		if err := r.Set(name, val); err != nil {
			return nil, err
		}
	}

	var err error
	if e.maxConcurrent > 0 {
		err = e.runAsync(ctx, r)
	} else {
		e.numInProgress.Add(1)
		err = r.RunContext(ctx)
		e.numInProgress.Add(-1)
	}

	// check for error and return variable
	if err != nil {
		return nil, err
	}
	if len(outVars) > 0 {
		results := []*tengo.Variable{}
		for _, name := range outVars {
			results = append(results, r.Get(name))
		}
		return results, nil
	}
	return r.GetAll(), nil
}

func (e *executor) runAsync(ctx context.Context, compiled Runnable) error {
	if e.sem == nil {
		return errors.New("bug: pool/limiter not defined")
	}
	if err := e.sem.Acquire(ctx, 1); err != nil {
		return err
	}
	e.numInProgress.Add(1)
	errCh := make(chan error)
	go func() {
		defer func() {
			e.numInProgress.Add(-1)
			defer e.sem.Release(1)
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
