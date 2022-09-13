package tnglib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtmlTemplate(t *testing.T) {
	err := runTengoScript("_testdata/htmltpl.tengo", "io", "html/template")
	assert.NoError(t, err)
}

func BenchmarkHtmlTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := runTengoScript("_testdata/htmltpl.tengo", "io", "html/template")
		if err != nil {
			b.Log(err)
		}
	}
}
