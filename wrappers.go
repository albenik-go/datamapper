package datamapper

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	V *time.Time
}

func (t *NullTime) Scan(value interface{}) error {
	if value == nil {
		*t.V = time.Time{}
		return nil
	}
	return convertAssign(t.V, value)
}

func (t *NullTime) Value() (driver.Value, error) {
	if t.V.IsZero() {
		return nil, nil
	}
	return *t.V, nil
}

type NullString struct {
	V *string
}

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s.V = ""
		return nil
	}
	return convertAssign(s.V, value)
}

func (s *NullString) Value() (driver.Value, error) {
	if *s.V == "" {
		return nil, nil
	}
	return *s.V, nil
}

type IntBool struct {
	V *bool
}

func (b *IntBool) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var i int
	if err := convertAssign(&i, value); err != nil {
		return err
	}
	*b.V = i != 0
	return nil
}

func (b *IntBool) Value() (driver.Value, error) {
	if *b.V {
		return 1, nil
	}
	return 0, nil
}
