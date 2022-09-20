package email_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	scriptFile := "../_testdata/message.tengo"
	err := script.RunFile(scriptFile, "fmt", "email", "context")
	assert.NoError(t, err)
}
