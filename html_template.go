package tnglib

import (
	"html/template"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

var htmlTplModule = map[string]tengo.Object{
	// new(string) => template
	"new": &tengo.UserFunction{
		Name:  "new",
		Value: htmlTplST(template.New),
	},
	// parse(name, string) => template
	"parse": &tengo.UserFunction{
		Name:  "parse",
		Value: htmlTplParse,
	},
	// execute_string(name, string, data) => string
	"execute_string": &tengo.UserFunction{
		Name:  "execute_string",
		Value: htmlTplParseExec,
	},
	// parse_files(...string) => template
	"parse_files": &tengo.UserFunction{
		Name:  "parse_files",
		Value: htmlTplASTE(template.ParseFiles),
	},
	// parse_glob(string) => template
	"parse_glob": &tengo.UserFunction{
		Name:  "parse_glob",
		Value: htmlTplSTE(template.ParseGlob),
	},
}

func makeHtmlTemplate(tpl *template.Template) *tengo.ImmutableMap {
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
				Value: htmlTplTE(tpl.Clone),
			},
			// defined_templates
			"defined_templates": &tengo.UserFunction{
				Name:  "defined_templates",
				Value: stdlib.FuncARS(tpl.DefinedTemplates),
			},
			// new(string) => template
			"new": &tengo.UserFunction{
				Name:  "new",
				Value: htmlTplST(tpl.New),
			},
			// lookup(string) => template
			"lookup": &tengo.UserFunction{
				Name:  "lookup",
				Value: htmlTplST(tpl.Lookup),
			},
			// option(...string) => template
			"option": &tengo.UserFunction{
				Name:  "option",
				Value: htmlTplAST(tpl.Option),
			},
			// parse(string) => template
			"parse": &tengo.UserFunction{
				Name:  "parse",
				Value: htmlTplSTE(tpl.Parse),
			},
			// parse_files(names...) => template
			"parse_files": &tengo.UserFunction{
				Name:  "parse_files",
				Value: htmlTplASTE(tpl.ParseFiles),
			},
			// parse_glob(pattern) => template
			"parse_glob": &tengo.UserFunction{
				Name:  "parse_glob",
				Value: htmlTplSTE(tpl.ParseGlob),
			},
			// execute(writer, data) => error
			"execute": &tengo.UserFunction{
				Name:  "execute",
				Value: FuncAWARE(tpl.Execute),
			},
			// execute_string(writer, data) => error
			"execute_string": &tengo.UserFunction{
				Name:  "execute_string",
				Value: FuncAWAREs(tpl.Execute),
			},
		},
	}
}

func htmlTplParse(args ...tengo.Object) (tengo.Object, error) {
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
	return makeHtmlTemplate(tpl), nil
}
func htmlTplParseExec(args ...tengo.Object) (tengo.Object, error) {
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
	s := sb.String()
	if len(s) > tengo.MaxStringLen {
		return nil, tengo.ErrStringLimit
	}
	return &tengo.String{Value: s}, nil
}

func htmlTplTE(fn func() (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		t, err := fn()
		if err != nil {
			return WrapError(err), nil
		}
		return makeHtmlTemplate(t), nil
	}
}

func htmlTplST(fn func(string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		name, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		return makeHtmlTemplate(fn(name)), nil
	}
}
func htmlTplAST(fn func(...string) *template.Template) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		params, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		return makeHtmlTemplate(fn(params...)), nil
	}
}
func htmlTplSTE(fn func(string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		pattern, err := ArgToString(args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(pattern)
		if err != nil {
			return WrapError(err), nil
		}
		return makeHtmlTemplate(t), nil
	}
}

func htmlTplASTE(fn func(...string) (*template.Template, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		filenames, err := ArgsToStrings(1, args...)
		if err != nil {
			return nil, err
		}
		t, err := fn(filenames...)
		if err != nil {
			return WrapError(err), nil
		}
		return makeHtmlTemplate(t), nil
	}
}
