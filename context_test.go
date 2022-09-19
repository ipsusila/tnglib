package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	err := tnglib.RunTengoScriptFile("_testdata/context.tengo", "context")
	assert.NoError(t, err)
}
