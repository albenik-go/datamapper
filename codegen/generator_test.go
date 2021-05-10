package codegen_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/albenik-go/datamapper/codegen"
)

const generated = `// Generated code! DO NOT EDIT!.
// github.com/albenik-go/datamapper/cmd/dmgen (NO VERSION INFO)

package mapper

import (
	"github.com/albenik-go/datamapper"
)

// ModelMapperBase shared column list is always the same for all mapper instances.
var ModelMapperBase = struct {
	SelectColumns []string
	InsertColumns []string
	UpdateColumns []string
}{
	SelectColumns: []string{"id", "string", "bool", "wrapped_bool", "time"},
	InsertColumns: []string{"string", "bool", "wrapped_bool", "time"},
	UpdateColumns: []string{"string", "bool", "wrapped_bool", "time"},
}

var ModelModel = struct {
	ID          string
	String      string
	Bool        string
	WrappedBool string
	Time        string
}{
	ID:          "id",
	String:      "string",
	Bool:        "bool",
	WrappedBool: "wrapped_bool",
	Time:        "time",
}

type ModelWrapper struct {
	entity *Model
}

func (m *ModelWrapper) ID() datamapper.Field {
	return datamapper.Field{Name: "id", Ref: &m.entity.ID}
}

func (m *ModelWrapper) String() datamapper.Field {
	return datamapper.Field{Name: "string", Ref: &m.entity.String}
}

func (m *ModelWrapper) Bool() datamapper.Field {
	return datamapper.Field{Name: "bool", Ref: &m.entity.Bool}
}

func (m *ModelWrapper) WrappedBool() datamapper.Field {
	return datamapper.Field{Name: "wrapped_bool", Ref: &datamapper.IntBool{V: &m.entity.WrappedBool}}
}

func (m *ModelWrapper) Time() datamapper.Field {
	return datamapper.Field{Name: "time", Ref: &m.entity.Time}
}

type ModelMapper struct {
	entity *Model
	fields *ModelWrapper

	selectFields []interface{}
	insertFields []interface{}
	updateFields []interface{}
}

func NewModelMapper(e *Model) *ModelMapper {
	return &ModelMapper{
		entity: e,
		fields: &ModelWrapper{entity: e},

		selectFields: []interface{}{&e.ID, &e.String, &e.Bool, &datamapper.IntBool{V: &e.WrappedBool}, &e.Time},
		insertFields: []interface{}{&e.String, &e.Bool, &datamapper.IntBool{V: &e.WrappedBool}, &e.Time},
		updateFields: []interface{}{&e.String, &e.Bool, &datamapper.IntBool{V: &e.WrappedBool}, &e.Time},
	}
}

func (m *ModelMapper) EmptyClone() *ModelMapper {
	return NewModelMapper(new(Model))
}

func (m *ModelMapper) EmptyCloneAsInterface() interface{} {
	return m.EmptyClone()
}

func (m *ModelMapper) SelectColumns() []string {
	return ModelMapperBase.SelectColumns
}

func (m *ModelMapper) SelectFields() []interface{} {
	return m.selectFields
}

func (m *ModelMapper) InsertColumns() []string {
	return ModelMapperBase.InsertColumns
}

func (m *ModelMapper) InsertFields() []interface{} {
	return m.insertFields
}

func (m *ModelMapper) UpdateColumns() []string {
	return ModelMapperBase.UpdateColumns
}

func (m *ModelMapper) UpdateFields() []interface{} {
	return m.updateFields
}

func (m *ModelMapper) UpdateFieldsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           &m.entity.ID,
		"string":       &m.entity.String,
		"bool":         &m.entity.Bool,
		"wrapped_bool": &m.entity.WrappedBool,
		"time":         &m.entity.Time,
	}
}

func (m *ModelMapper) Model() *ModelWrapper {
	return m.fields
}

func (m *ModelMapper) Entity() *Model {
	return m.entity
}
`

func TestGenerate(t *testing.T) {
	t.Run("Defaults", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		err := codegen.SimplifiedGenerate(
			"../internal/testmodels/model.go",
			"mapper",
			"db",
			"db",
			nil,
			false,
			buf,
		)
		require.NoError(t, err)
		require.Equal(t, generated, buf.String())
	})
}
