//+build ignore

package testmodels

import (
	"time"
)

type Model struct {
	ID          int64     `col:"id,auto"`
	String      string    `col:"string"`
	Bool        bool      `col:"bool"`
	WrappedBool bool      `col:"wrapped_bool,wrap=IntBool"`
	Time        time.Time `col:"time"`
}
