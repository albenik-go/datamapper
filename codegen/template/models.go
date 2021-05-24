package template

type Header struct {
	Package      string
	DmgenVersion string
}

type ModelInfo struct {
	EntityType string

	SelectFields []*FieldInfo // All fields
	InsertFields []*FieldInfo
	UpdateFields []*FieldInfo
}

type FieldInfo struct {
	FieldName  string
	ColumnName string
	Wrapper    string
}
