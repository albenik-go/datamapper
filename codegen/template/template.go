package template

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode/utf8"
)

func parseTemplate(s string) *template.Template {
	t := template.New("mapper").
		Funcs(template.FuncMap{
			"lcFirst": func(s string) string {
				r, size := utf8.DecodeRuneInString(s)
				return strings.ToLower(string(r)) + s[size:]
			},
			"asColumnsSlice": func(fields []*FieldInfo) string {
				cols := make([]string, len(fields))
				for i, f := range fields {
					cols[i] = fmt.Sprintf("%q", f.ColumnName)
				}
				return fmt.Sprintf("{%s}", strings.Join(cols, ", "))
			},
			"asRefsSlice": func(fields []*FieldInfo) string {
				refs := make([]string, len(fields))
				for i, f := range fields {
					if f.Nullable {
						refs[i] = fmt.Sprintf("&datamapper.Nullable{WrappedValue: &m.%s}", f.FieldName)
					} else {
						refs[i] = fmt.Sprintf("&m.%s", f.FieldName)
					}
				}
				return fmt.Sprintf("{%s}", strings.Join(refs, ", "))
			},
		})

	return template.Must(t.Parse(s))
}

func WriteHeader(w io.Writer, i *Header) error {
	return headerTpl.Execute(w, i)
}

func WriteModel(w io.Writer, m *ModelInfo) error {
	return bodyTpl.Execute(w, m)
}
