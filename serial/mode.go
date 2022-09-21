package serial

import (
	"errors"
	"strconv"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	ser "go.bug.st/serial"
)

var (
	mParity = map[string]ser.Parity{
		"N": ser.NoParity,
		"O": ser.OddParity,
		"E": ser.EvenParity,
		"M": ser.MarkParity,
		"S": ser.SpaceParity,
	}
	mStopBits = map[string]ser.StopBits{
		"1":   ser.OneStopBit,
		"1.5": ser.OnePointFiveStopBits,
		"2":   ser.TwoStopBits,
	}
)

// first argument could be
// 1. EMPTY (NONE)
// 2. int - baud rate
// 3. string - 9600,8,n,1
func newMode() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) > 1 {
			return nil, tengo.ErrWrongNumArguments
		}

		// create serial.Mode instance
		m := defaultMode()
		if len(args) == 1 {
			if baud, ok := tengo.ToInt(args[0]); ok {
				m.BaudRate = baud
			} else if short, ok := tengo.ToString(args[0]); ok {
				if err := assignMode(&m, short); err != nil {
					return tnglib.WrapError(err), nil
				}
			}
		}
		return modeToObject(m), nil
	}
}

func defaultMode() ser.Mode {
	return ser.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   ser.NoParity,
		StopBits: ser.OneStopBit,
		InitialStatusBits: &ser.ModemOutputBits{
			RTS: true,
			DTR: true,
		},
	}
}

func assignMode(m *ser.Mode, short string) error {
	// e.g. 9600-8-N-1 OR 9600,8,N,1
	items := strings.FieldsFunc(short, func(sep rune) bool {
		return sep == '-' || sep == ',' || sep == '/'
	})
	if len(items) != 4 {
		return errors.New("invalid serial config: " + short)
	}

	// 1. Baudrate
	baud, err := strconv.Atoi(items[0])
	if err != nil {
		return err
	}
	m.BaudRate = baud

	// 2. Databits
	dataBits, err := strconv.Atoi(items[1])
	if err != nil {
		return err
	}
	m.DataBits = dataBits

	// 3. Parity
	var ok bool
	parity := strings.ToUpper(items[2])
	m.Parity, ok = mParity[parity]
	if !ok {
		return errors.New("invalid parity: " + parity)
	}

	// 4. StopBits
	stopBits := items[3]
	m.StopBits, ok = mStopBits[stopBits]
	if !ok {
		return errors.New("invalid stop bits: " + items[3])
	}
	return nil
}

func modeToObject(m ser.Mode) *tengo.Map {
	// default value
	if m.InitialStatusBits == nil {
		m.InitialStatusBits = &ser.ModemOutputBits{
			RTS: true,
			DTR: true,
		}
	}
	return &tengo.Map{
		Value: map[string]tengo.Object{
			"baud_rate": &tengo.Int{Value: int64(m.BaudRate)},
			"data_bits": &tengo.Int{Value: int64(m.DataBits)},
			"parity":    &tengo.Int{Value: int64(m.Parity)},
			"stop_bits": &tengo.Int{Value: int64(m.StopBits)},
			"initial_status_bits": &tengo.Map{
				Value: map[string]tengo.Object{
					"rts": tnglib.Ternary(m.InitialStatusBits.RTS, tengo.TrueValue, tengo.FalseValue),
					"dtr": tnglib.Ternary(m.InitialStatusBits.DTR, tengo.TrueValue, tengo.FalseValue),
				},
			},
		},
	}
}
func objectToMode(obj tengo.Object) ser.Mode {
	m := defaultMode()
	m.BaudRate = tnglib.MapGet(obj, "baud_rate", m.BaudRate, tengo.ToInt)
	m.DataBits = tnglib.MapGet(obj, "data_bits", m.DataBits, tengo.ToInt)
	m.Parity = ser.Parity(tnglib.MapGet(obj, "parity", int(m.Parity), tengo.ToInt))
	m.StopBits = ser.StopBits(tnglib.MapGet(obj, "stop_bits", int(m.StopBits), tengo.ToInt))

	vo := tnglib.MapGet(obj, "initial_status_bits", nil, tnglib.ToObject)
	m.InitialStatusBits = &ser.ModemOutputBits{
		RTS: tnglib.MapGet(vo, "rts", true, tengo.ToBool),
		DTR: tnglib.MapGet(vo, "dtr", true, tengo.ToBool),
	}

	return m
}
