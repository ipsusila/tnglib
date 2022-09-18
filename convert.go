package tnglib

import (
	"fmt"
	"time"

	"github.com/d5/tengo/v2"
)

// ArgToString convert tengo function call arguments to string.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgToString(args ...tengo.Object) (string, error) {
	if len(args) != 1 {
		return "", tengo.ErrWrongNumArguments
	}
	str, ok := tengo.ToString(args[0])
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	return str, nil
}

// ArgIToString convert tengo function call arguments to string.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToString(idx int, args ...tengo.Object) (string, error) {
	if idx >= len(args) {
		return "", tengo.ErrWrongNumArguments
	}
	str, ok := tengo.ToString(args[idx])
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "string(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return str, nil
}

// ArgIToInt convert tengo function call arguments to int.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToInt(idx int, args ...tengo.Object) (int, error) {
	if idx >= len(args) {
		return 0, tengo.ErrWrongNumArguments
	}
	n, ok := tengo.ToInt(args[idx])
	if !ok {
		return 0, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "integer(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return n, nil
}

// ArgIToBool convert tengo function call arguments to boolean.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToBool(idx int, args ...tengo.Object) (bool, error) {
	if idx >= len(args) {
		return false, tengo.ErrWrongNumArguments
	}
	v, ok := tengo.ToBool(args[idx])
	if !ok {
		return false, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "bool(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return v, nil
}

// ArgIToInt64 convert tengo function call arguments to string.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToInt64(idx int, args ...tengo.Object) (int64, error) {
	if idx >= len(args) {
		return 0, tengo.ErrWrongNumArguments
	}
	v, ok := tengo.ToInt64(args[idx])
	if !ok {
		return 0, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "integer(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return v, nil
}

// ArgToByteSlice convert tengo function call arguments to []byte.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgToByteSlice(args ...tengo.Object) ([]byte, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}
	data, ok := tengo.ToByteSlice(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	return data, nil
}

// ArgIToByteSlice convert tengo function call arguments to []byte.
// If the argument count is not equal to one, it will return ErrWrongNumArguments
func ArgIToByteSlice(idx int, args ...tengo.Object) ([]byte, error) {
	if idx >= len(args) {
		return nil, tengo.ErrWrongNumArguments
	}
	data, ok := tengo.ToByteSlice(args[idx])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "bytes(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return data, nil
}

// ArgsToStrings convert tengo function call arguments to string slice.
// If the argument count is less than minArg, it will return ErrWrongNumArguments
func ArgsToStrings(minArg int, args ...tengo.Object) ([]string, error) {
	if len(args) < minArg {
		return nil, tengo.ErrWrongNumArguments
	}
	items := []string{}
	for idx, arg := range args {
		filename, ok := tengo.ToString(arg)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("arg[%d]", idx),
				Expected: "string(compatible)",
				Found:    arg.TypeName(),
			}
		}
		items = append(items, filename)
	}
	return items, nil
}

// ArgIToContext convert argument to context value
func ArgIToContext(idx int, args ...tengo.Object) (*Context, error) {
	if idx >= len(args) {
		return nil, tengo.ErrWrongNumArguments
	}
	// get context
	ctx, ok := args[0].(*Context)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "context",
			Found:    args[idx].TypeName(),
		}
	}
	return ctx, nil
}

// ArgToContext convert argument to context value
func ArgToContext(args ...tengo.Object) (*Context, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}
	// get context
	ctx, ok := args[0].(*Context)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "context",
			Found:    args[0].TypeName(),
		}
	}
	return ctx, nil
}

// ArgIToTime convert argument to context value
func ArgIToTime(idx int, args ...tengo.Object) (time.Time, error) {
	if idx >= len(args) {
		return time.Time{}, tengo.ErrWrongNumArguments
	}
	// get time
	tm, ok := tengo.ToTime(args[idx])
	if !ok {
		return time.Time{}, tengo.ErrInvalidArgumentType{
			Name:     fmt.Sprintf("arg[%d]", idx),
			Expected: "time(compatible)",
			Found:    args[idx].TypeName(),
		}
	}
	return tm, nil
}

// ObjectMapToMap converts map[string]tengo.Object to map[string]interface{}
func ObjectMapToMap(v map[string]tengo.Object) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range v {
		m[key] = tengo.ToInterface(val)
	}
	return m
}

// MapToObject convert map[string]interface{} to tengo.Object
func MapToObject(v map[string]interface{}) (tengo.Object, error) {
	kv := make(map[string]tengo.Object)
	for vk, vv := range v {
		vo, err := tengo.FromInterface(vv)
		if err != nil {
			return nil, err
		}
		kv[vk] = vo
	}
	return &tengo.Map{Value: kv}, nil
}

// ArgIToMap convert argument to context value
func ArgIToMap(idx int, args ...tengo.Object) (map[string]interface{}, error) {
	if idx >= len(args) {
		return nil, tengo.ErrWrongNumArguments
	}

	switch m := args[idx].(type) {
	case *tengo.Map:
		return ObjectMapToMap(m.Value), nil
	case *tengo.ImmutableMap:
		return ObjectMapToMap(m.Value), nil
	}

	return nil, tengo.ErrInvalidArgumentType{
		Name:     fmt.Sprintf("arg[%d]", idx),
		Expected: "map(compatible)",
		Found:    args[idx].TypeName(),
	}
}
