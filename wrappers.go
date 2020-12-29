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

type IntBool struct {
	WrappedValue *bool
}

func (b *IntBool) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var i int
	if err := convertAssign(&i, value); err != nil {
		return err
	}
	*b.WrappedValue = i != 0
	return nil
}

func (b *IntBool) Value() (driver.Value, error) {
	if *b.WrappedValue {
		return 1, nil
	}
	return 0, nil
}
