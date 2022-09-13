package tnglib

import (
	"strings"
	"text/template"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

var txtTplModule = map[string]tengo.Object{
	// new(string) => template
	"new": &tengo.UserFunction{
		Name:  "new",
		Value: txtTplST(template.New),
	},
	// parse(name, string) => template
	"parse": &tengo.UserFunction{
		Name:  "parse",
		Value: txtTplParse,
	},
	// execute_string(name, string, data) => string
	"execute_string": &tengo.UserFunction{
		Name:  "execute_string",
		Value: txtTplParseExec,
	},
	// parse_files(...string) => template
	"parse_files": &tengo.UserFunction{
		Name:  "parse_files",
		Value: txtTplASTE(template.ParseFiles),
	},
	// parse_glob(string) => template
	"parse_glob": &tengo.UserFunction{
		Name:  "parse_glob",
		Value: txtTplSTE(template.ParseGlob),
	},
}

func makeTxtTemplate(tpl *template.Template) *tengo.ImmutableMap {
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
				Value: txtTplTE(tpl.Clone),
			},
			// defined_templates
			"defined_templates": &tengo.UserFunction{
				Name:  "defined_templates",
				Value: stdlib.FuncARS(tpl.DefinedTemplates),
			},
			// new(string) => template
			"new": &tengo.UserFunction{
				Name:  "new",
				Value: txtTplST(tpl.New),
			},
			// lookup(string) => template
			"lookup": &tengo.UserFunction{
				Name:  "lookup",
				Value: txtTplST(tpl.Lookup),
			},
			// option(...string) => template
			"option": &tengo.UserFunction{
				Name:  "option",
				Value: txtTplAST(tpl.Option),
			},
			// parse(string) => template
			"parse": &tengo.UserFunction{
				Name:  "parse",
				Value: txtTplSTE(tpl.Parse),
			},
			// parse_files(names...) => template
			"parse_files": &tengo.UserFunction{
				Name:  "parse_files",
				Value: txtTplASTE(tpl.ParseFiles),
			},
			// parse_glob(pattern) => template
			"parse_glob": &tengo.UserFunction{
				Name:  "parse_glob",
				Value: txtTplSTE(tpl.ParseGlob),
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

func txtTplParse(args ...tengo.Object) (tengo.Object, error) {
	name, err := ArgIToString(0, args...)
	if err != nil {
		return nil, err
	}
	text, err := ArgIToString(1, args...)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(name).Parse(text)
	if err != nil {
		return wrapError(err), nil
	}
	return makeTxtTemplate(tpl), nil
}
func txtTplParseExec(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 3 {
		return nil, tengo.ErrWrongNumArguments
	}
	name, err := ArgIToString(0, args...)
	if err != nil {
		return nil, err
	}
	text, err := ArgIToString(1, args...)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New(name).Parse(text)
	if err != nil {
		return wrapError(err), nil
	}

	var sb strings.Builder
	data := tengo.ToInterface(args[2])
	if err := tpl.Execute(&sb, data); err != nil {
		return wrapError(err), nil
	}
	return &tengo.String{Value: sb.String()}, nil
}

func txtTplTE(fn func() (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		t, err := fn()
		if err != nil {
			return wrapError(err), nil
		}
		return makeTxtTemplate(t), nil
	}
}

func txtTplST(fn func(string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		name, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		return makeTxtTemplate(fn(name)), nil
	}
}
func txtTplAST(fn func(...string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		params, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		return makeTxtTemplate(fn(params...)), nil
	}
}
func txtTplSTE(fn func(string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		pattern, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(pattern)
		if err != nil {
			return wrapError(err), nil
		}
		return makeTxtTemplate(t), nil
	}
}

func txtTplASTE(fn func(...string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		filenames, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(filenames...)
		if err != nil {
			return wrapError(err), nil
		}
		return makeTxtTemplate(t), nil
	}
}
