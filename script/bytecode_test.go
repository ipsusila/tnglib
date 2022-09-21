package script_test

import (
	"context"
	"testing"
	"time"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestBytecodeRead(t *testing.T) {
	by, err := script.BytecodeFromCompiled("../_testdata/work.out")
	assert.NoError(t, err)
	t.Log("Compiled at: ", by.CompiledAt().Format(time.RFC3339))
	t.Log("Config: ", by.Configuration())
	t.Log("Bytecode: ", by.String())

	r := by.Runnable()
	ctx, cancel := context.WithTimeout(context.TODO(), 4*time.Second)
	defer cancel()
	err = r.RunContext(ctx)
	assert.NoError(t, err)
}
