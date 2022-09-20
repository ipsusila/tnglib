package email

import (
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Msg struct {
	tengo.ImmutableMap
	Msg *mail.Email
}

func (m *Msg) TypeName() string {
	return "email.msg"
}
func (m *Msg) String() string {
	return fmt.Sprintf("<email.msg> ch: %s", m.Msg.Charset)
}
func (m *Msg) Copy() tengo.Object {
	return &Msg{Msg: m.Msg, ImmutableMap: m.ImmutableMap}
}
func (m *Msg) IsFalsy() bool {
	return m.Msg == nil
}
func (m *Msg) Equals(x tengo.Object) bool {
	if x == nil || m == x {
		return m == x
	}
	v, ok := x.(*Msg)
	if !ok {
		return false
	}
	return v.Msg == m.Msg
}

func makeMsgMap(m *Msg) tengo.ImmutableMap {
	msg := m.Msg
	return tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"encoding": &tengo.Int{Value: int64(msg.Encoding)},
			"charset":  &tengo.String{Value: msg.Charset},
			"dkim_msg": &tengo.String{Value: msg.DkimMsg},
			"get_error": &tengo.UserFunction{
				Name:  "get_error",
				Value: stdlib.FuncARE(msg.GetError),
			},
			"get_from": &tengo.UserFunction{
				Name:  "get_from",
				Value: stdlib.FuncARS(msg.GetFrom),
			},
			"get_message": &tengo.UserFunction{
				Name:  "get_message",
				Value: stdlib.FuncARS(msg.GetMessage),
			},
			"get_recipients": &tengo.UserFunction{
				Name:  "get_recipients",
				Value: stdlib.FuncARSs(msg.GetRecipients),
			},
			"add_addresses": &tengo.UserFunction{
				Name:  "add_addresses",
				Value: funcASSsRM(m, msg.AddAddresses),
			},
			"add_header": &tengo.UserFunction{
				Name:  "add_header",
				Value: funcASSsRM(m, msg.AddHeader),
			},
			"add_bcc": &tengo.UserFunction{
				Name:  "add_bcc",
				Value: funcASsRM(m, msg.AddBcc),
			},
			"add_cc": &tengo.UserFunction{
				Name:  "add_cc",
				Value: funcASsRM(m, msg.AddCc),
			},
			"add_to": &tengo.UserFunction{
				Name:  "add_to",
				Value: funcASsRM(m, msg.AddTo),
			},
			"set_date": &tengo.UserFunction{
				Name:  "set_date",
				Value: funcASRM(m, msg.SetDate),
			},
			"set_from": &tengo.UserFunction{
				Name:  "set_from",
				Value: funcASRM(m, msg.SetFrom),
			},
			"set_list_unsubscribe": &tengo.UserFunction{
				Name:  "set_list_unsubscribe",
				Value: funcASRM(m, msg.SetListUnsubscribe),
			},
			"set_reply_to": &tengo.UserFunction{
				Name:  "set_reply_to",
				Value: funcASRM(m, msg.SetReplyTo),
			},
			"set_return_path": &tengo.UserFunction{
				Name:  "set_return_path",
				Value: funcASRM(m, msg.SetReturnPath),
			},
			"set_sender": &tengo.UserFunction{
				Name:  "set_sender",
				Value: funcASRM(m, msg.SetSender),
			},
			"set_subject": &tengo.UserFunction{
				Name:  "set_subject",
				Value: funcASRM(m, msg.SetSubject),
			},
			"add_alternative": &tengo.UserFunction{
				Name:  "add_alternative",
				Value: funcAddAlternative(m),
			},
			"add_alternative_data": &tengo.UserFunction{
				Name:  "add_alternative_data",
				Value: funcAddAlternativeData(m),
			},
			"set_body": &tengo.UserFunction{
				Name:  "set_body",
				Value: funcSetBody(m),
			},
			"set_body_data": &tengo.UserFunction{
				Name:  "set_body_data",
				Value: funcSetBodyData(m),
			},
			"attach": &tengo.UserFunction{
				Name:  "attach",
				Value: funcAttach(m),
			},
			"set_priority": &tengo.UserFunction{
				Name:  "set_priority",
				Value: funcSetPriority(m),
			},
			"set_dkim": &tengo.UserFunction{
				Name:  "set_dkim",
				Value: funcSetDkim(m),
			},
		},
	}
}

/*
TODO:
func (email *Email) AddHeaders(headers textproto.MIMEHeader) *Email
*/

// NewMessage create new email message
func NewMessage(m *Msg) *Msg {
	if m == nil {
		m = &Msg{Msg: mail.NewMSG()}
		im := makeMsgMap(m)
		m.ImmutableMap = im
	}
	return m
}

func newMsgFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		return NewMessage(nil), nil
	}
}

func funcASSsRM(msg *Msg, fn func(string, ...string) *mail.Email) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) < 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, err := tnglib.ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		argsV := []string{}
		if len(args) > 1 {
			argsV, err = tnglib.ArgsToStrings(0, args[:1]...)
			if err != nil {
				return nil, err
			}
		}
		fn(s1, argsV...)
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
func funcASsRM(msg *Msg, fn func(...string) *mail.Email) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		argsV, err := tnglib.ArgsToStrings(0, args...)
		if err != nil {
			return nil, err
		}
		fn(argsV...)
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
func funcASRM(msg *Msg, fn func(string) *mail.Email) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s1, err := tnglib.ArgToString(args...)
		if err != nil {
			return nil, err
		}
		fn(s1)
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
func funcAddAlternative(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		i1, err := tnglib.ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		body, err := tnglib.ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}
		switch i1 {
		case int(mail.TextPlain):
			msg.Msg.AddAlternative(mail.TextPlain, body)
		case int(mail.TextHTML):
			msg.Msg.AddAlternative(mail.TextHTML, body)
		case int(mail.TextCalendar):
			msg.Msg.AddAlternative(mail.TextCalendar, body)
		}
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}

		return msg, nil
	}
}
func funcAddAlternativeData(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		i1, err := tnglib.ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		data, err := tnglib.ArgIToByteSlice(1, args...)
		if err != nil {
			return nil, err
		}
		switch i1 {
		case int(mail.TextPlain):
			msg.Msg.AddAlternativeData(mail.TextPlain, data)
		case int(mail.TextHTML):
			msg.Msg.AddAlternativeData(mail.TextHTML, data)
		case int(mail.TextCalendar):
			msg.Msg.AddAlternativeData(mail.TextCalendar, data)
		}
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}

		return msg, nil
	}
}
func funcSetBody(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		i1, err := tnglib.ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		body, err := tnglib.ArgIToString(1, args...)
		if err != nil {
			return nil, err
		}
		switch i1 {
		case int(mail.TextPlain):
			msg.Msg.SetBody(mail.TextPlain, body)
		case int(mail.TextHTML):
			msg.Msg.SetBody(mail.TextHTML, body)
		case int(mail.TextCalendar):
			msg.Msg.SetBody(mail.TextCalendar, body)
		}
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}

		return msg, nil
	}
}
func funcSetBodyData(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		i1, err := tnglib.ArgIToInt(0, args...)
		if err != nil {
			return nil, err
		}
		data, err := tnglib.ArgIToByteSlice(1, args...)
		if err != nil {
			return nil, err
		}
		switch i1 {
		case int(mail.TextPlain):
			msg.Msg.SetBodyData(mail.TextPlain, data)
		case int(mail.TextHTML):
			msg.Msg.SetBodyData(mail.TextHTML, data)
		case int(mail.TextCalendar):
			msg.Msg.SetBodyData(mail.TextCalendar, data)
		}
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}

		return msg, nil
	}
}

func funcAttach(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		fi := objectToFile(args[0])
		msg.Msg.Attach(&fi)
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
func funcSetPriority(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		i1, err := tnglib.ArgToInt(args...)
		if err != nil {
			return nil, err
		}
		if i1 == int(mail.PriorityHigh) {
			msg.Msg.SetPriority(mail.PriorityHigh)
		} else {
			msg.Msg.SetPriority(mail.PriorityLow)
		}
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
func funcSetDkim(msg *Msg) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		so := objectToDkimOptions(args[0])
		msg.Msg.SetDkim(so)
		if err := msg.Msg.GetError(); err != nil {
			return tnglib.WrapError(err), nil
		}
		return msg, nil
	}
}
