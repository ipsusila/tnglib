package tnglib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	err := runTengoScript("_testdata/context.tengo", "context")
	assert.NoError(t, err)
}
