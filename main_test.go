package main //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTarget(t *testing.T) {
	require := require.New(t)
	t.Parallel()

	tables := []struct {
		name     string
		t        string
		expected Target
		err      error
	}{
		{
			"001",
			"foobar.tld",
			Target{"foobar.tld", 5201},
			nil,
		},
		{
			"002",
			"foobar.tld:1234",
			Target{"foobar.tld", 1234},
			nil,
		},
		{
			"003",
			"foobar:foobar:foobar",
			Target{},
			ErrCouldNotDetermineTarget,
		},
	}

	for _, table := range tables {
		table := table
		t.Run(table.name, func(t *testing.T) {
			t.Parallel()
			trgt, err := NewTarget(table.t)
			require.ErrorIs(err, table.err)
			require.Equal(table.expected, trgt)
		})
	}
}
