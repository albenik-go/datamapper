package template

const header = `// Generated code! DO NOT EDIT!.
// dmgen {{.DmgenVersion}}

package {{.Package}}

import (
	"github.com/albenik-go/datamapper"
)`

const body = `

// Columns list is always the same for all mapper instances.
// So let's keep it pre-created and re-use it.
var {{.ModelName}}MapperBase = struct {
	SelectColumns []string
	InsertColumns []string
	UpdateColumns []string
}{
	SelectColumns: []string{{.SelectFields | asColumnsSlice}},
	InsertColumns: []string{{.InsertFields | asColumnsSlice}},
	UpdateColumns: []string{{.UpdateFields | asColumnsSlice}},
}

var {{.ModelName}}Model = struct {
{{- range .SelectFields}}
	{{.FieldName}} string
{{- end}}
}{
{{- range .SelectFields}}
	{{.FieldName}}: "{{.ColumnName}}",
{{- end}}
}

type {{.ModelName}}Wrapper struct {
	entity *{{.ModelName}}
}
{{range .SelectFields}}
func (m *{{$.ModelName}}Wrapper) {{.FieldName}}() datamapper.Field {
	return datamapper.Field{Name: "{{.ColumnName}}", Ref: {{. | asRef}}}
}
{{end}}

type {{.ModelName}}Mapper struct {
	entity *{{.ModelName}}
	fields *{{.ModelName}}Wrapper

	selectFields []interface{}
	insertFields []interface{}
	updateFields []interface{}
}

func New{{.ModelName}}Mapper(m *{{.ModelName}}) *{{.ModelName}}Mapper {
	return &{{.ModelName}}Mapper{
		entity:  m,
		fields: &{{.ModelName}}Wrapper{entity: m},

		selectFields: []interface{}{{.SelectFields | asRefsSlice}},
		insertFields: []interface{}{{.InsertFields | asRefsSlice}},
		updateFields: []interface{}{{.UpdateFields | asRefsSlice}},
	}
}

func (m *{{.ModelName}}Mapper) SelectColumns() []string {
	return {{.ModelName}}MapperBase.SelectColumns
}

func (m *{{.ModelName}}Mapper) SelectFields() []interface{} {
	return m.selectFields
}

func (m *{{.ModelName}}Mapper) InsertColumns() []string {
	return {{.ModelName}}MapperBase.InsertColumns
}

func (m *{{.ModelName}}Mapper) InsertFields() []interface{} {
	return m.insertFields
}

func (m *{{.ModelName}}Mapper) UpdateColumns() []string {
	return {{.ModelName}}MapperBase.UpdateColumns
}

func (m *{{.ModelName}}Mapper) UpdateFields() []interface{} {
	return m.updateFields
}

func (m *{{.ModelName}}Mapper) Model() *{{.ModelName}}Wrapper {
	return m.fields
}

func (m *{{.ModelName}}Mapper) Entity() *{{.ModelName}} {
	return m.entity
}`

var (
	headerTpl = parseTemplate(header)
	bodyTpl   = parseTemplate(body)
)
