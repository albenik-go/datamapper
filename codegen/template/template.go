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
			"asRef":          asRefFilter,
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

func asRefsSliceFilter(receiver string, fields []*FieldInfo) string {
	refs := make([]string, len(fields))
	for i, f := range fields {
		s := fmt.Sprintf("%s.%s", receiver, f.FieldName)
		if len(f.Wrapper) > 0 {
			s = fmt.Sprintf("&datamapper.%s{V: %s}", f.Wrapper, s)
		}
		refs[i] = s
	}
	return fmt.Sprintf("{%s}", strings.Join(refs, ", "))
}

func asRefFilter(receiver string, field *FieldInfo) string {
	s := fmt.Sprintf("%s.entity.%s", receiver, field.FieldName)
	if len(field.Wrapper) > 0 {
		s = fmt.Sprintf("&datamapper.%s{V: %s}", field.Wrapper, s)
	}
	return s
}

func WriteHeader(w io.Writer, i *Header) error {
	return headerTpl.Execute(w, i)
}

func WriteModel(w io.Writer, m *ModelInfo) error {
	return bodyTpl.Execute(w, m)
}
