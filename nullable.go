package datamapper

import (
	"database/sql/driver"
)

type Nullable struct {
	WrappedValue interface{}
}

func (n *Nullable) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return convertAssign(&n.WrappedValue, value)
}

func (n *Nullable) Value() (driver.Value, error) {
	return n.WrappedValue, nil
}
