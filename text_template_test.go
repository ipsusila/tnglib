package tnglib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextTemplate(t *testing.T) {
	err := runTengoScript("_testdata/texttpl.tengo", "io", "text/template")
	assert.NoError(t, err)
}

func BenchmarkTextTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := runTengoScript("_testdata/texttpl.tengo", "io", "text/template")
		if err != nil {
			b.Log(err)
		}
	}
}
