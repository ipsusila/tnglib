package email

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

// modules name
var (
	Name = "email"
)

var (
	// Module registered here
	emailModule = map[string]tengo.Object{
		// constants
		"encoding_none":             &tengo.Int{Value: int64(mail.EncodingNone)},
		"encoding_base64":           &tengo.Int{Value: int64(mail.EncodingBase64)},
		"encoding_quoted_printable": &tengo.Int{Value: int64(mail.EncodingQuotedPrintable)},
		"text_plain":                &tengo.Int{Value: int64(mail.TextPlain)},
		"text_html":                 &tengo.Int{Value: int64(mail.TextHTML)},
		"text_calendar":             &tengo.Int{Value: int64(mail.TextCalendar)},
		"priority_low":              &tengo.Int{Value: int64(mail.PriorityLow)},
		"priority_high":             &tengo.Int{Value: int64(mail.PriorityHigh)},

		// smtp_conn(string, keep_alive, timeout) => cli
		"smtp_connect": &tengo.UserFunction{
			Name:  "smtp_connect",
			Value: smtpConnectFunc(),
		},
		// new_msg()
		"new_msg": &tengo.UserFunction{
			Name:  "new_msg",
			Value: newMsgFunc(),
		},
		// new_file() => mail
		"new_file": &tengo.UserFunction{
			Name:  "new_file",
			Value: newFile(),
		},
		// new_dkim_sig_options() => dkim.SigOptions
		"new_dkim_sig_options": &tengo.UserFunction{
			Name:  "new_dkim_sig_options",
			Value: newDkimSigOptions(),
		},
	}
)

func init() {
	// register module
	tnglib.RegisterModule(Name, emailModule)
}
