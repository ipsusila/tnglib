package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib"
	"github.com/stretchr/testify/assert"
)

func TestIoScanner(t *testing.T) {
	err := tnglib.RunTengoScriptFile("_testdata/scanner.tengo", "io")
	assert.NoError(t, err)
}
