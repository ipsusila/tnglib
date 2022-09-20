package email

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	"github.com/toorop/go-dkim"
)

func makeDkimSig(dk *dkim.SigOptions) *tengo.Map {
	if dk == nil {
		so := dkim.NewSigOptions()
		dk = &so
	}
	return &tengo.Map{
		Value: map[string]tengo.Object{
			"version":                 &tengo.Int{Value: int64(dk.Version)},
			"private_key":             &tengo.Bytes{Value: dk.PrivateKey},
			"domain":                  &tengo.String{Value: dk.Domain},
			"selector":                &tengo.String{Value: dk.Selector},
			"auid":                    &tengo.String{Value: dk.Auid},
			"canonicalization":        &tengo.String{Value: dk.Canonicalization},
			"algo":                    &tengo.String{Value: dk.Algo},
			"headers":                 tnglib.StringsToObject(dk.Headers),
			"body_length":             &tengo.Int{Value: int64(dk.BodyLength)},
			"query_methods":           tnglib.StringsToObject(dk.QueryMethods),
			"add_signature_timestamp": tnglib.BoolObject(dk.AddSignatureTimestamp),
			"signature_expire_in":     &tengo.Int{Value: int64(dk.SignatureExpireIn)},
			"copied_header_fields":    tnglib.StringsToObject(dk.CopiedHeaderFields),
		},
	}
}

func objectToDkimOptions(obj tengo.Object) dkim.SigOptions {
	op := dkim.NewSigOptions()
	if obj == nil {
		return op
	}
	var mv map[string]tengo.Object
	switch v := obj.(type) {
	case *tengo.Map:
		mv = v.Value
	case *tengo.ImmutableMap:
		mv = v.Value
	}
	if len(mv) == 0 {
		return op
	}

	for key, val := range mv {
		switch key {
		case "version":
			if v, ok := tengo.ToInt64(val); ok {
				op.Version = uint(v)
			}
		case "private_key":
			if v, ok := tengo.ToByteSlice(val); ok {
				op.PrivateKey = v
			}
		case "domain":
			if v, ok := tengo.ToString(val); ok {
				op.Domain = v
			}
		case "selector":
			if v, ok := tengo.ToString(val); ok {
				op.Selector = v
			}
		case "auid":
			if v, ok := tengo.ToString(val); ok {
				op.Auid = v
			}
		case "canonicalization":
			if v, ok := tengo.ToString(val); ok {
				op.Canonicalization = v
			}
		case "algo":
			if v, ok := tengo.ToString(val); ok {
				op.Algo = v
			}
		case "headers":
			if va, err := tnglib.ObjectToStrings(val); err == nil {
				op.Headers = va
			}
		case "body_length":
			if v, ok := tengo.ToInt64(val); ok {
				op.BodyLength = uint(v)
			}
		case "query_methods":
			if va, err := tnglib.ObjectToStrings(val); err == nil {
				op.QueryMethods = va
			}
		case "add_signature_timestamp":
			if v, ok := tengo.ToBool(val); ok {
				op.AddSignatureTimestamp = v
			}
		case "signature_expire_in":
			if v, ok := tengo.ToInt64(val); ok {
				op.SignatureExpireIn = uint64(v)
			}
		case "copied_header_fields":
			if va, err := tnglib.ObjectToStrings(val); err == nil {
				op.CopiedHeaderFields = va
			}
		}
	}

	return op
}

func newDkimSigOptions() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		return makeDkimSig(nil), nil
	}
}
