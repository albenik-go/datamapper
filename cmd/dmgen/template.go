package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode/utf8"
)

const raw = `import (
	"github.com/albenik-go/datamapper"
)

var {{.ModelName | lcFirst}}MapperBase = &datamapper.ModelMapperBase{
	SelectColumns: []string{{"{"}}{{.SelectFields | asColumnsSlice}}{{"}"}},
	InsertColumns: []string{{"{"}}{{.InsertFields | asColumnsSlice}}{{"}"}},
	UpdateColumns: []string{{"{"}}{{.UpdateFields | asColumnsSlice}}{{"}"}},
}

type {{.ModelName}}Wrapper struct {
	model *{{.ModelName}}
}
{{range .SelectFields}}
func (m *{{$.ModelName}}Wrapper) {{.FieldName}}() datamapper.Field {
	return datamapper.Field{Name: "{{.ColumnName}}", Ref: &m.model.{{.FieldName}}}
}
{{end}}

type {{.ModelName}}Mapper struct {
	model  *{{.ModelName}}
	base   *datamapper.ModelMapperBase
	fields *{{.ModelName}}Wrapper

	selectFields []interface{}
	insertFields []interface{}
	updateFields []interface{}
}

func New{{.ModelName}}Mapper(m *{{.ModelName}}) *{{.ModelName}}Mapper {
	return &{{.ModelName}}Mapper{
		model:  m,
		base:   modelMapperBase,
		fields: &{{.ModelName}}Wrapper{model: m},

		selectFields: []interface{}{{"{"}}{{.SelectFields | asRefsSlice}}{{"}"}},
		insertFields: []interface{}{{"{"}}{{.InsertFields | asRefsSlice}}{{"}"}},
		updateFields: []interface{}{{"{"}}{{.UpdateFields | asRefsSlice}}{{"}"}},
	}
}

func (m *{{.ModelName}}Mapper) SelectColumns() []string {
	return m.base.SelectColumns
}

func (m *{{.ModelName}}Mapper) SelectFields() []interface{} {
	return m.selectFields
}

func (m *{{.ModelName}}Mapper) InsertColumns() []string {
	return m.base.InsertColumns
}

func (m *{{.ModelName}}Mapper) InsertFields() []interface{} {
	return m.insertFields
}

func (m *{{.ModelName}}Mapper) UpdateColumns() []string {
	return m.base.UpdateColumns
}

func (m *{{.ModelName}}Mapper) UpdateFields() []interface{} {
	return m.updateFields
}

func (m *{{.ModelName}}Mapper) Model() *{{.ModelName}}Wrapper {
	return m.fields
}`

var tpl = parseTemplate()

func parseTemplate() *template.Template {
	t := template.New("mapper").
		Funcs(template.FuncMap{
			"lcFirst": func(s string) string {
				r, size := utf8.DecodeRuneInString(s)
				return strings.ToLower(string(r)) + s[size:]
			},
			"asColumnsSlice": func(fields []*fieldInfo) string {
				cols := make([]string, len(fields))
				for i, f := range fields {
					cols[i] = fmt.Sprintf("%q", f.ColumnName)
				}
				return strings.Join(cols, ", ")
			},
			"asRefsSlice": func(fields []*fieldInfo) string {
				refs := make([]string, len(fields))
				for i, f := range fields {
					refs[i] = fmt.Sprintf("&m.%s", f.FieldName)
				}
				return strings.Join(refs, ", ")
			},
		})

	return template.Must(t.Parse(raw))
}

func execTemplate(w io.Writer, m *modelInfo) {
	if err := tpl.Execute(w, m); err != nil {
		panic(err)
	}
}

type fieldInfo struct {
	FieldName  string
	ColumnName string
}

type modelInfo struct {
	ModelName string

	SelectFields []*fieldInfo // All fields
	InsertFields []*fieldInfo
	UpdateFields []*fieldInfo
}
