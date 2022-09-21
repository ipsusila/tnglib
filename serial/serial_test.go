package serial_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestSerial(t *testing.T) {
	err := script.RunFile("../_testdata/serial.tengo", "fmt", "serial")
	assert.NoError(t, err)
}
