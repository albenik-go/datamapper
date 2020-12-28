package datamapper

type Nullable struct {
	Value interface{}
}

// Scan implements the Scanner interface.
func (n *Nullable) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return convertAssign(&n.Value, value)
}
