package datamapper

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	V *time.Time
}

func (w *NullTime) Scan(value interface{}) error {
	if value == nil {
		*w.V = time.Time{}
		return nil
	}
	return ConvertAssign(w.V, value)
}

func (w *NullTime) Value() (driver.Value, error) {
	if w.V.IsZero() {
		return nil, nil
	}
	return *w.V, nil
}

type NullInt struct {
	V *int
}

func (w *NullInt) Scan(value interface{}) error {
	if value == nil {
		*w.V = 0
		return nil
	}
	return ConvertAssign(w.V, value)
}

func (w *NullInt) Value() (driver.Value, error) {
	if *w.V == 0 {
		return nil, nil
	}
	return *w.V, nil
}

type NullUint64 struct {
	V *uint64
}

func (w *NullUint64) Scan(value interface{}) error {
	if value == nil {
		*w.V = 0
		return nil
	}
	return ConvertAssign(w.V, value)
}

func (w *NullUint64) Value() (driver.Value, error) {
	if *w.V == 0 {
		return nil, nil
	}
	return *w.V, nil
}

type NullString struct {
	V *string
}

func (w *NullString) Scan(value interface{}) error {
	if value == nil {
		*w.V = ""
		return nil
	}
	return ConvertAssign(w.V, value)
}

func (w *NullString) Value() (driver.Value, error) {
	if *w.V == "" {
		return nil, nil
	}
	return *w.V, nil
}

type IntBool struct {
	V *bool
}

func (w *IntBool) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var i int
	if err := ConvertAssign(&i, value); err != nil {
		return err
	}
	*w.V = i != 0

	return nil
}

func (w *IntBool) Value() (driver.Value, error) {
	if *w.V {
		return int64(1), nil
	}
	return int64(0), nil
}
