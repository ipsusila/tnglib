package serial

import (
	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/ipsusila/tnglib"
	ser "go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func makePort(p ser.Port) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"read": &tengo.UserFunction{
				Name:  "read",
				Value: stdlib.FuncAYRIE(p.Read),
			},
			"write": &tengo.UserFunction{
				Name:  "write",
				Value: stdlib.FuncAYRIE(p.Write),
			},
			"close": &tengo.UserFunction{
				Name:  "close",
				Value: stdlib.FuncARE(p.Close),
			},
			"reset_input_buffer": &tengo.UserFunction{
				Name:  "reset_input_buffer",
				Value: stdlib.FuncARE(p.ResetInputBuffer),
			},
			"reset_output_buffer": &tengo.UserFunction{
				Name:  "reset_output_buffer",
				Value: stdlib.FuncARE(p.ResetOutputBuffer),
			},
			"set_rts": &tengo.UserFunction{
				Name:  "set_rts",
				Value: tnglib.FuncABRE(p.SetRTS),
			},
			"set_dtr": &tengo.UserFunction{
				Name:  "set_dtr",
				Value: tnglib.FuncABRE(p.SetDTR),
			},
			"set_read_timeout": &tengo.UserFunction{
				Name:  "set_read_timeout",
				Value: tnglib.FuncADRE(p.SetReadTimeout),
			},
		},
	}
}

func openPort() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}
		portName, err := tnglib.ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}

		m := defaultMode()

		// second argument is short config, e.g. 9600-N-8-1
		short, ok := tengo.ToString(args[1])
		if ok {
			if err := assignMode(&m, short); err != nil {
				return tnglib.WrapError(err), nil
			}
		} else {
			m = objectToMode(args[1])
		}

		com, err := ser.Open(portName, &m)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makePort(com), nil
	}
}

func portDetailToObject(p *enumerator.PortDetails) tengo.Object {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			"name":          &tengo.String{Value: p.Name},
			"is_usb":        tnglib.BoolObject(p.IsUSB),
			"vid":           &tengo.String{Value: p.VID},
			"pid":           &tengo.String{Value: p.PID},
			"serial_number": &tengo.String{Value: p.SerialNumber},
			"product":       &tengo.String{Value: p.Product},
		},
	}
}

func enumDetailedPorts() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			return tnglib.WrapError(err), nil
		}

		va := tengo.Array{}
		for _, p := range ports {
			va.Value = append(va.Value, portDetailToObject(p))
		}
		return &va, nil
	}
}
