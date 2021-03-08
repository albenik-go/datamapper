package datamapper

import (
	"testing"
	"time"
)

func TestNullable_Scan(t *testing.T) {
	var (
		t1 = time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC)
		t2 time.Time
	)
	if err := convertAssign(&NullTime{V: &t2}, t1); err != nil {
		t.Fatal(err)
	}
	if !t2.Equal(t1) {
		t.Fatal("Time is not equal!", t1, "vs", t2)
	}
}
