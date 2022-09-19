package tnglib_test

import (
	"testing"

	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestHtmlTemplate(t *testing.T) {
	err := script.RunFile("_testdata/htmltpl.tengo", "fmt", "os", "io", "html/template")
	assert.NoError(t, err)
}

func BenchmarkHtmlTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := script.RunFile("_testdata/htmltpl.tengo", "fmt", "os", "io", "html/template")
		if err != nil {
			b.Log(err)
		}
	}
}
