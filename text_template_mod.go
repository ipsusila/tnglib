package tnglib

import (
	"text/template"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

func makeTxtTemplate(tpl *template.Template) *tengo.ImmutableMap {
	return &tengo.ImmutableMap{
		Value: map[string]tengo.Object{
			// name() => string
			"name": &tengo.UserFunction{
				Name:  "name",
				Value: stdlib.FuncARS(tpl.Name),
			},
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
			"execute": &tengo.UserFunction{
				Name:  "execute",
				Value: tplExecute(tpl.Execute),
			},
		},
	}
}
