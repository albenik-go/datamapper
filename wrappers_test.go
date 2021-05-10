package datamapper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/albenik-go/datamapper"
)

func TestNullable_Scan(t *testing.T) {
	var (
		time1 = time.Date(2020, time.December, 31, 23, 59, 59, 0, time.UTC)
		time2 time.Time
	)

	require.NoError(t, datamapper.ConvertAssign(&datamapper.NullTime{V: &time2}, time1))
	require.Equal(t, time2, time1)
}
