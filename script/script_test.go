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
	conf := script.DefaultConfig()
	conf.Modules = []string{"fmt", "times", "context"}
	man := script.NewManager()
	exe := script.NewExecutor(man, maxConcurrent)
	err := man.Add(id, "../_testdata/work.tengo", conf)
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
}

func TestUnlimited(t *testing.T) {
	execScript(t, script.Unlimited, 10, 4*time.Second)
}

func TestLimited(t *testing.T) {
	execScript(t, 4, 10, 15*time.Second)
}
