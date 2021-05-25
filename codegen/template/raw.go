package template

const header = `// Generated code! DO NOT EDIT!.
// github.com/albenik-go/datamapper/cmd/dmgen {{.DmgenVersion}}

package {{.Package}}

import (
	"github.com/albenik-go/datamapper"
)`

const body = `

// {{.EntityType}}MapperBase shared column list is always the same for all mapper instances.
var {{.EntityType}}MapperBase = struct {
	SelectColumns []string
	InsertColumns []string
	UpdateColumns []string
}{
	SelectColumns: []string{{.SelectFields | asColumnsSlice}},
	InsertColumns: []string{{.InsertFields | asColumnsSlice}},
	UpdateColumns: []string{{.UpdateFields | asColumnsSlice}},
}

var {{.EntityType}}Model = struct {
{{- range .SelectFields}}
	{{.FieldName}} string
{{- end}}
}{
{{- range .SelectFields}}
	{{.FieldName}}: "{{.ColumnName}}",
{{- end}}
}

type {{.EntityType}}EntityWrapper struct {
	entity *{{.EntityType}}
}
{{range .SelectFields}}
func (m *{{$.EntityType}}EntityWrapper) {{.FieldName}}() datamapper.Field {
	return datamapper.Field{Name: "{{.ColumnName}}", Ref: {{. | asRef "&m"}}}
}
{{end}}

type {{.EntityType}}Mapper struct {
	entity *{{.EntityType}}
	fields *{{.EntityType}}EntityWrapper

	selectFields []interface{}
	insertFields []interface{}
	updateFields []interface{}
}

func New{{.EntityType}}Mapper(e *{{.EntityType}}) *{{.EntityType}}Mapper {
	if e == nil {
		e = new({{.EntityType}})
	}

	return &{{.EntityType}}Mapper{
		entity:  e,
		fields: &{{.EntityType}}EntityWrapper{entity: e},

		selectFields: []interface{}{{.SelectFields | asRefsSlice "&e"}},
		insertFields: []interface{}{{.InsertFields | asRefsSlice "&e"}},
		updateFields: []interface{}{{.UpdateFields | asRefsSlice "&e"}},
	}
}

func (m *{{.EntityType}}Mapper) SelectColumns() []string {
	return {{.EntityType}}MapperBase.SelectColumns
}

func (m *{{.EntityType}}Mapper) SelectFields() []interface{} {
	return m.selectFields
}

func (m *{{.EntityType}}Mapper) InsertColumns() []string {
	return {{.EntityType}}MapperBase.InsertColumns
}

func (m *{{.EntityType}}Mapper) InsertFields() []interface{} {
	return m.insertFields
}

func (m *{{.EntityType}}Mapper) UpdateColumns() []string {
	return {{.EntityType}}MapperBase.UpdateColumns
}

func (m *{{.EntityType}}Mapper) UpdateFields() []interface{} {
	return m.updateFields
}

func (m *{{.EntityType}}Mapper) UpdateFieldsMap() map[string]interface{} {
	return map[string]interface{}{
		{{- range .UpdateFields}}
			"{{.ColumnName}}": {{. | asRef "&m"}},
		{{- end}}
	}
}

func (m *{{.EntityType}}Mapper) Model() *{{.EntityType}}EntityWrapper {
	return m.fields
}

func (m *{{.EntityType}}Mapper) Entity() *{{.EntityType}} {
	return m.entity
}

func (m *{{.EntityType}}Mapper) EmptyClone() *{{.EntityType}}Mapper {
	return New{{.EntityType}}Mapper(new({{.EntityType}}))
}

func (m *{{.EntityType}}Mapper) UntypedEntity() interface{} {
	return m.entity
}

func (m *{{.EntityType}}Mapper) UntypedEmptyClone() interface{} {
	return m.EmptyClone()
}
{{if .AutoincrementField}}
func (m *{{.EntityType}}Mapper) SetID(id {{.AutoincrementField.FieldType}}) {
	m.entity.{{.AutoincrementField.FieldName}} = id
}{{end}}
`

var (
	headerTpl = parseTemplate(header)
	bodyTpl   = parseTemplate(body)
)
