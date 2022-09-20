package email

import (
	"path/filepath"

	"github.com/d5/tengo/v2"
	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

func makeFile(fi *mail.File) *tengo.Map {
	return &tengo.Map{
		Value: map[string]tengo.Object{
			"file_path": &tengo.String{Value: fi.FilePath},
			"name":      &tengo.String{Value: fi.Name},
			"mime_type": &tengo.String{Value: fi.MimeType},
			"b64data":   &tengo.String{Value: fi.B64Data},
			"data":      &tengo.Bytes{Value: fi.Data},
			"inline":    tnglib.BoolObject(fi.Inline),
		},
	}
}

func objectToFile(obj tengo.Object) mail.File {
	fi := mail.File{}
	if obj == nil {
		return fi
	}
	var mv map[string]tengo.Object
	switch v := obj.(type) {
	case *tengo.Map:
		mv = v.Value
	case *tengo.ImmutableMap:
		mv = v.Value
	}
	if len(mv) == 0 {
		return fi
	}

	for key, val := range mv {
		switch key {
		case "file_path":
			fi.FilePath, _ = tengo.ToString(val)
			if fi.FilePath != "" {
				fi.Name = filepath.Base(fi.FilePath)
			}
		case "name":
			fi.Name, _ = tengo.ToString(val)
		case "mime_type":
			fi.MimeType, _ = tengo.ToString(val)
		case "b64data":
			fi.B64Data, _ = tengo.ToString(val)
		case "data":
			fi.Data, _ = tengo.ToByteSlice(val)
		case "inline":
			fi.Inline, _ = tengo.ToBool(val)
		}
	}

	return fi
}

func newFile() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		var err error
		fi := &mail.File{}
		if len(args) > 0 {
			fi.FilePath, err = tnglib.ArgIToString(0, args...)
			if err != nil {
				return nil, err
			}
			fi.Name = filepath.Base(fi.FilePath)
		}
		if len(args) > 1 {
			name, err := tnglib.ArgIToString(1, args...)
			if err != nil {
				return nil, err
			}
			fi.Name = name
		}
		if len(args) > 2 {
			mime, err := tnglib.ArgIToString(2, args...)
			if err != nil {
				return nil, err
			}
			fi.MimeType = mime
		}
		if len(args) > 3 {
			b64, err := tnglib.ArgIToString(3, args...)
			if err != nil {
				return nil, err
			}
			fi.B64Data = b64
		}
		if len(args) > 4 {
			data, err := tnglib.ArgIToByteSlice(4, args...)
			if err != nil {
				return nil, err
			}
			fi.Data = data
		}
		if len(args) > 5 {
			inline, err := tnglib.ArgIToBool(5, args...)
			if err != nil {
				return nil, err
			}
			fi.Inline = inline
		}
		if len(args) > 6 {
			return nil, tengo.ErrWrongNumArguments
		}

		return makeFile(fi), nil
	}
}
