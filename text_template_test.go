package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestTextTemplate(t *testing.T) {
	err := script.RunFile("_testdata/texttpl.tengo", "fmt", "os", "io", "text/template")
	assert.NoError(t, err)
}

func BenchmarkTextTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := script.RunFile("_testdata/texttpl.tengo", "fmt", "os", "io", "text/template")
		if err != nil {
			b.Log(err)
		}
	}
}
