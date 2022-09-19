package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	err := script.RunFile("_testdata/context.tengo", "fmt", "times", "context")
	assert.NoError(t, err)
}
