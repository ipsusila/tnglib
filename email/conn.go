package email

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

func makeConnection(conn *mail.SMTPClient) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"keep_alive":   tnglib.BoolObject(conn.KeepAlive),
			"send_timeout": &tengo.Int{Value: int64(conn.SendTimeout)},
			"close": &tengo.UserFunction{
				Name:  "close",
				Value: stdlib.FuncARE(conn.Close),
			},
			"reset": &tengo.UserFunction{
				Name:  "reset",
				Value: stdlib.FuncARE(conn.Reset),
			},
			"noop": &tengo.UserFunction{
				Name:  "noop",
				Value: stdlib.FuncARE(conn.Noop),
			},
			"quit": &tengo.UserFunction{
				Name:  "quit",
				Value: stdlib.FuncARE(conn.Quit),
			},
			"send": &tengo.UserFunction{
				Name:  "send",
				Value: sendMessage(conn),
			},
			"send_envelope_from": &tengo.UserFunction{
				Name:  "send_envelope_from",
				Value: sendEnvelopeFrom(conn),
			},
		},
	}
}

// handle 2 cases:
// 1. func SendMessage(from string, recipients []string, msg string, client *SMTPClient) error
// 2. func Email.Send(client *SMTPClient) error
func sendMessage(conn *mail.SMTPClient) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 1 {
			msg, ok := args[0].(*Msg)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "email.msg",
					Found:    args[0].TypeName(),
				}
			}
			err := msg.Msg.Send(conn)
			return tnglib.WrapError(err), nil
		} else if len(args) == 3 {
			from, err := tnglib.ArgIToString(0, args...)
			if err != nil {
				return nil, err
			}

			to, err := tnglib.ObjectToStrings(args[1])
			if err != nil {
				return nil, err
			}
			msg, err := tnglib.ArgIToString(2, args...)
			if err != nil {
				return nil, err
			}
			err = mail.SendMessage(from, to, msg, conn)
			return tnglib.WrapError(err), nil
		}

		return nil, tengo.ErrWrongNumArguments
	}
}

func sendEnvelopeFrom(conn *mail.SMTPClient) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		from, err := tnglib.ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		msg, ok := args[1].(*Msg)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "email.msg",
				Found:    args[1].TypeName(),
			}
		}
		err = msg.Msg.SendEnvelopeFrom(from, conn)
		return tnglib.WrapError(err), nil
	}
}
