package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func New(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Funcs(functions).Parse(tmpl))
}

func formatField(field string) (string, error) {
	values := strings.Split(field, ":")
	return strings.Join(values[:2], " "), nil
}

func formatFieldOptional(field string) (string, error) {
	values := strings.Split(field, ":")
	if len(values) < 2 {
		return "", fmt.Errorf("invalid format for field: %s", field)
	}
	return values[0] + " *" + values[1], nil
}

func formatFieldName(field string) string {
	values := strings.Split(field, ":")
	return values[0]
}

func formatFieldTag(field string) (string, error) {
	values := strings.Split(field, ":")
	if len(values) < 2 {
		return "", fmt.Errorf("invalid format for field: %s", field)
	}
	return strings.ReplaceAll(strcase.ToLowerCamel(values[0]), "ID", "Id"), nil
}

func formatValidation(field string) string {
	values := strings.Split(field, ":")
	if len(values) < 3 {
		return "-"
	}
	return values[2]
}

func hasValidation(field string) bool {
	return len(strings.Split(field, ":")) > 2
}

var functions = template.FuncMap{
	"formatField":         formatField,
	"formatFieldTag":      formatFieldTag,
	"formatValidation":    formatValidation,
	"hasValidation":       hasValidation,
	"formatFieldName":     formatFieldName,
	"formatFieldOptional": formatFieldOptional,
	"toCamel":             strcase.ToCamel,
}
