package template

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

func parseTemplate(s string) *template.Template {
	t := template.New("mapper").
		Funcs(template.FuncMap{
			"asColumnsSlice": asColumnSliceFilter,
			"asRefsSlice":    asRefsSliceFilter,
		})

	return template.Must(t.Parse(s))
}

func asColumnSliceFilter(fields []*FieldInfo) string {
	cols := make([]string, len(fields))
	for i, f := range fields {
		cols[i] = fmt.Sprintf("%q", f.ColumnName)
	}
	return fmt.Sprintf("{%s}", strings.Join(cols, ", "))
}

func asRefsSliceFilter(fields []*FieldInfo) string {
	refs := make([]string, len(fields))
	for i, f := range fields {
		rs := fmt.Sprintf("&m.%s", f.FieldName)
		if len(f.Wrappers) > 0 {
			for i := len(f.Wrappers) - 1; i >= 0; i-- {
				rs = fmt.Sprintf("&datamapper.%s{WrappedValue: %s}", f.Wrappers[i], rs)
			}
		}
		refs[i] = rs
	}
	return fmt.Sprintf("{%s}", strings.Join(refs, ", "))
}

func WriteHeader(w io.Writer, i *Header) error {
	return headerTpl.Execute(w, i)
}

func WriteModel(w io.Writer, m *ModelInfo) error {
	return bodyTpl.Execute(w, m)
}
