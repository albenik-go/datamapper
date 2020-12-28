package template

const header = `// Generated code! DO NOT EDIT!.

package {{.Package}}

import (
	"github.com/albenik-go/datamapper"
)`

const body = `

// Columns list is always the same for all mapper instances.
// So let's keep it pre-created and re-use it.
var {{.ModelName | lcFirst}}MapperBase = [][]string {
	{{.SelectFields | asColumnsSlice}},
	{{.InsertFields | asColumnsSlice}},
	{{.UpdateFields | asColumnsSlice}},
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
	base   [][]string
	fields *{{.ModelName}}Wrapper

	selectFields []interface{}
	insertFields []interface{}
	updateFields []interface{}
}

func New{{.ModelName}}Mapper(m *{{.ModelName}}) *{{.ModelName}}Mapper {
	return &{{.ModelName}}Mapper{
		model:  m,
		base:   {{.ModelName | lcFirst}}MapperBase,
		fields: &{{.ModelName}}Wrapper{model: m},

		selectFields: []interface{}{{.SelectFields | asRefsSlice}},
		insertFields: []interface{}{{.InsertFields | asRefsSlice}},
		updateFields: []interface{}{{.UpdateFields | asRefsSlice}},
	}
}

func (m *{{.ModelName}}Mapper) SelectColumns() []string {
	return m.base[0]
}

func (m *{{.ModelName}}Mapper) SelectFields() []interface{} {
	return m.selectFields
}

func (m *{{.ModelName}}Mapper) InsertColumns() []string {
	return m.base[1]
}

func (m *{{.ModelName}}Mapper) InsertFields() []interface{} {
	return m.insertFields
}

func (m *{{.ModelName}}Mapper) UpdateColumns() []string {
	return m.base[2]
}

func (m *{{.ModelName}}Mapper) UpdateFields() []interface{} {
	return m.updateFields
}

func (m *{{.ModelName}}Mapper) Model() *{{.ModelName}}Wrapper {
	return m.fields
}`

var (
	headerTpl = parseTemplate(header)
	bodyTpl   = parseTemplate(body)
)
