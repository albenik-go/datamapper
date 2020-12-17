package datamapper

type ModelMapperBase struct {
	SelectColumns []string
	InsertColumns []string
	UpdateColumns []string
}

type Field struct {
	Name string
	Ref  interface{}
}

type Fieldset []*Field

func NewFieldset(fields ...*Field) Fieldset {
	return fields
}

func (s Fieldset) SelectColumns() []string {
	names := make([]string, len(s))
	for i, f := range s {
		names[i] = f.Name
	}
	return names
}

func (s Fieldset) SelectFields() []interface{} {
	refs := make([]interface{}, len(s))
	for i, f := range s {
		refs[i] = f.Ref
	}
	return refs
}
