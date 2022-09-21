package script_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func execScript(t *testing.T, maxConcurrent, n int, timeout time.Duration) {
	id := "test"
	modules := []string{"fmt", "times", "context"}
	scriptFile := "../_testdata/work.tengo"
	compiledFile := "../_testdata/work.out"

	conf := script.DefaultConfig()
	conf.InitVars = map[string]interface{}{
		"X":       100,
		"message": "hello world",
	}
	conf.Modules = modules
	man := script.NewManager()
	exe := script.NewExecutor(man, maxConcurrent)
	err := man.AddFile(id, scriptFile, conf)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, err := exe.ExecContext(ctx, id, nil)
			assert.NoError(t, err)
		}()
		time.Sleep(1 * time.Millisecond)
		t.Log("Num in progress: ", exe.NumInProgress())
	}
	wg.Wait()

	// save entry to file
	entry := man.Entry(id)
	err = entry.SaveTo(compiledFile)
	assert.NoError(t, err)

	// run entry
	err = script.RunFile(compiledFile, modules...)
	assert.NoError(t, err)
}

func TestUnlimited(t *testing.T) {
	execScript(t, script.Unlimited, 10, 4*time.Second)
}

func TestLimited(t *testing.T) {
	execScript(t, 4, 10, 15*time.Second)
}
