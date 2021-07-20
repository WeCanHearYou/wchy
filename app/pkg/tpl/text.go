package tpl

import (
	"github.com/getfider/fider/app/pkg/markdown"
	"strings"
	"text/template"
)

func GetTextTemplate(name string, rawText string) (*template.Template, error) {
	tpl, err := template.New(name).Funcs(templateFunctions).Funcs(template.FuncMap{
		"markdown": func(input string) string {
			return markdown.PlainText(input)
		},
	}).Parse(rawText)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func Execute(tmpl *template.Template, data interface{}) (string, error) {
	builder := &strings.Builder{}
	if err := tmpl.Execute(builder, data); err != nil {
		return "", err
	}
	return builder.String(), nil
}
