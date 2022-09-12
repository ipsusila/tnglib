package tnglib

import (
	"fmt"
	"text/template"

	"github.com/d5/tengo/v2"
)

var txtTplModule = map[string]tengo.Object{
	"new": &tengo.UserFunction{
		Name:  "new",
		Value: txtTemplateNew,
	},
	"parse": &tengo.UserFunction{
		Name:  "parse",
		Value: txtTemplateParse,
	},
	"parse_files": &tengo.UserFunction{
		Name:  "parse_files",
		Value: txtTemplateParseFiles,
	},
	"parse_glob": &tengo.UserFunction{
		Name:  "parse_glob",
		Value: txtTemplateParseGlob,
	},
}

func txtTemplateNew(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}
	name, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	tpl := template.New(name)
	return makeTxtTemplate(tpl), nil
}

func txtTemplateParse(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}
	name, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	text, ok := tengo.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	tpl, err := template.New(name).Parse(text)
	if err != nil {
		return wrapError(err), nil
	}
	return makeTxtTemplate(tpl), nil
}

func txtTemplateParseFiles(args ...tengo.Object) (tengo.Object, error) {
	if len(args) < 1 {
		return nil, tengo.ErrWrongNumArguments
	}
	filenames := []string{}
	for idx, arg := range args {
		filename, ok := tengo.ToString(arg)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("arg%d", idx),
				Expected: "string(compatible)",
				Found:    arg.TypeName(),
			}
		}
		filenames = append(filenames, filename)
	}
	tpl, err := template.ParseFiles(filenames...)
	if err != nil {
		return wrapError(err), nil
	}
	return makeTxtTemplate(tpl), nil
}

func txtTemplateParseGlob(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}
	pattern, ok := tengo.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	tpl, err := template.ParseGlob(pattern)
	if err != nil {
		return wrapError(err), nil
	}
	return makeTxtTemplate(tpl), nil
}
