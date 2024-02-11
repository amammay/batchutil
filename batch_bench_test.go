package batchutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkAll(b *testing.B) {
	in := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		in = append(in, i)
	}

	b.ResetTimer()
	err := All(in, 1000, func([]int) error {
		return nil
	})
	require.NoError(b, err)
}
