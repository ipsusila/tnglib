package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestIoScanner(t *testing.T) {
	err := script.RunFile("_testdata/scanner.tengo", "fmt", "os", "io")
	assert.NoError(t, err)
}
