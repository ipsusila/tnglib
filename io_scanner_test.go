package tnglib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIoScanner(t *testing.T) {
	err := runTengoScript("_testdata/scanner.tengo", "io")
	assert.NoError(t, err)
}
