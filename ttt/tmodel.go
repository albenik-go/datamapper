package ttt

import (
	"time"
)

type Model struct {
	ID          int64     `db:"id,auto"`
	String      string    `db:"string"`
	Bool        bool      `db:"bool"`
	WrappedBool bool      `db:"wrapped_bool,wrap=IntBool"`
	Time        time.Time `db:"time"`
}
