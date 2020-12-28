package template

type Header struct {
	Package string
}

type ModelInfo struct {
	ModelName string

	SelectFields []*FieldInfo // All fields
	InsertFields []*FieldInfo
	UpdateFields []*FieldInfo
}

type FieldInfo struct {
	FieldName  string
	ColumnName string
	Nullable   bool
}
