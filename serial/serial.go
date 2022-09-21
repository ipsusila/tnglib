package serial

// https://github.com/bugst/go-serial

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	ser "go.bug.st/serial"
)

// modules name
var (
	Name = "serial"
)

var (
	// Module registered here
	serialModule = map[string]tengo.Object{
		"no_parity":                &tengo.Int{Value: int64(ser.NoParity)},
		"odd_parity":               &tengo.Int{Value: int64(ser.OddParity)},
		"even_parity":              &tengo.Int{Value: int64(ser.EvenParity)},
		"mark_parity":              &tengo.Int{Value: int64(ser.MarkParity)},
		"space_parity":             &tengo.Int{Value: int64(ser.SpaceParity)},
		"port_busy":                &tengo.Int{Value: int64(ser.PortBusy)},
		"port_not_found":           &tengo.Int{Value: int64(ser.PortNotFound)},
		"invalid_serial_port":      &tengo.Int{Value: int64(ser.InvalidSerialPort)},
		"permission_denied":        &tengo.Int{Value: int64(ser.PermissionDenied)},
		"invalid_speed":            &tengo.Int{Value: int64(ser.InvalidSpeed)},
		"invalid_data_bits":        &tengo.Int{Value: int64(ser.InvalidDataBits)},
		"invalid_parity":           &tengo.Int{Value: int64(ser.InvalidParity)},
		"invalid_stop_bits":        &tengo.Int{Value: int64(ser.InvalidStopBits)},
		"invalid_timeout_value":    &tengo.Int{Value: int64(ser.InvalidTimeoutValue)},
		"error_enumerating_ports":  &tengo.Int{Value: int64(ser.ErrorEnumeratingPorts)},
		"port_closed":              &tengo.Int{Value: int64(ser.PortClosed)},
		"function_not_implemented": &tengo.Int{Value: int64(ser.FunctionNotImplemented)},
		"one_stop_bit":             &tengo.Int{Value: int64(ser.OneStopBit)},
		"one_point_five_stop_bits": &tengo.Int{Value: int64(ser.OnePointFiveStopBits)},
		"two_stop_bits":            &tengo.Int{Value: int64(ser.TwoStopBits)},
		"get_ports_list": &tengo.UserFunction{
			Name:  "get_ports_list",
			Value: tnglib.FuncARSsE(ser.GetPortsList),
		},
		"mode": &tengo.UserFunction{
			Name:  "mode",
			Value: newMode(),
		},
		"open": &tengo.UserFunction{
			Name:  "open",
			Value: openPort(),
		},
		"get_detailed_ports_list": &tengo.UserFunction{
			Name:  "get_detailed_ports_list",
			Value: enumDetailedPorts(),
		},
	}
)

func init() {
	// register module
	tnglib.RegisterModule(Name, serialModule)
}
