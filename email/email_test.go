package email_test

import (
	"testing"

	"github.com/ipsusila/tnglib/email"
	"github.com/ipsusila/tnglib/script"
	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	conf, err := email.LoadConfig("../_testdata/private_conf.hjson")
	assert.NoError(t, err)

	err = email.RegisterSmtpServer("test", conf)
	assert.NoError(t, err)

	scriptFile := "../_testdata/email.tengo"
	err = script.RunFile(scriptFile, "fmt", "email", "os")
	assert.NoError(t, err)
}
