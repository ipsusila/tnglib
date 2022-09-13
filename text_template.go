package tnglib

import (
	"strings"
	"text/template"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

var textTplModule = map[string]tengo.Object{
	// new(string) => template
	"new": &tengo.UserFunction{
		Name:  "new",
		Value: textTplST(template.New),
	},
	// parse(name, string) => template
	"parse": &tengo.UserFunction{
		Name:  "parse",
		Value: textTplParse,
	},
	// execute_string(name, string, data) => string
	"execute_string": &tengo.UserFunction{
		Name:  "execute_string",
		Value: textTplParseExec,
	},
	// parse_files(...string) => template
	"parse_files": &tengo.UserFunction{
		Name:  "parse_files",
		Value: textTplASTE(template.ParseFiles),
	},
	// parse_glob(string) => template
	"parse_glob": &tengo.UserFunction{
		Name:  "parse_glob",
		Value: textTplSTE(template.ParseGlob),
	},
}

func makeTextTemplate(tpl *template.Template) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// name() => string
			"name": &tengo.UserFunction{
				Name:  "name",
				Value: stdlib.FuncARS(tpl.Name),
			},
			// clone() => template
			"clone": &tengo.UserFunction{
				Name:  "clone",
				Value: textTplTE(tpl.Clone),
			},
			// defined_templates
			"defined_templates": &tengo.UserFunction{
				Name:  "defined_templates",
				Value: stdlib.FuncARS(tpl.DefinedTemplates),
			},
			// new(string) => template
			"new": &tengo.UserFunction{
				Name:  "new",
				Value: textTplST(tpl.New),
			},
			// lookup(string) => template
			"lookup": &tengo.UserFunction{
				Name:  "lookup",
				Value: textTplST(tpl.Lookup),
			},
			// option(...string) => template
			"option": &tengo.UserFunction{
				Name:  "option",
				Value: textTplAST(tpl.Option),
			},
			// parse(string) => template
			"parse": &tengo.UserFunction{
				Name:  "parse",
				Value: textTplSTE(tpl.Parse),
			},
			// parse_files(names...) => template
			"parse_files": &tengo.UserFunction{
				Name:  "parse_files",
				Value: textTplASTE(tpl.ParseFiles),
			},
			// parse_glob(pattern) => template
			"parse_glob": &tengo.UserFunction{
				Name:  "parse_glob",
				Value: textTplSTE(tpl.ParseGlob),
			},
			// execute(writer, data) => error
			"execute": &tengo.UserFunction{
				Name:  "execute",
				Value: FuncWIE(tpl.Execute),
			},
			// execute_string(writer, data) => error
			"execute_string": &tengo.UserFunction{
				Name:  "execute_string",
				Value: FuncWISE(tpl.Execute),
			},
		},
	}
}

func textTplParse(args ...tengo.Object) (tengo.Object, error) {
	name, err := ArgIToString(0, args...)
	if err != nil {
		return nil, err
	}
	str, err := ArgIToString(1, args...)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(name).Parse(str)
	if err != nil {
		return WrapError(err), nil
	}
	return makeTextTemplate(tpl), nil
}
func textTplParseExec(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 3 {
		return nil, tengo.ErrWrongNumArguments
	}
	name, err := ArgIToString(0, args...)
	if err != nil {
		return nil, err
	}
	str, err := ArgIToString(1, args...)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(name).Parse(str)
	if err != nil {
		return WrapError(err), nil
	}

	var sb strings.Builder
	data := tengo.ToInterface(args[2])
	if err := tpl.Execute(&sb, data); err != nil {
		return WrapError(err), nil
	}
	return &tengo.String{Value: sb.String()}, nil
}

func textTplTE(fn func() (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		t, err := fn()
		if err != nil {
			return WrapError(err), nil
		}
		return makeTextTemplate(t), nil
	}
}

func textTplST(fn func(string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		name, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		return makeTextTemplate(fn(name)), nil
	}
}
func textTplAST(fn func(...string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		params, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		return makeTextTemplate(fn(params...)), nil
	}
}
func textTplSTE(fn func(string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		pattern, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(pattern)
		if err != nil {
			return WrapError(err), nil
		}
		return makeTextTemplate(t), nil
	}
}

func textTplASTE(fn func(...string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		filenames, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(filenames...)
		if err != nil {
			return WrapError(err), nil
		}
		return makeTextTemplate(t), nil
	}
}
