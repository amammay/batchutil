package batchutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBatchChanTimeout(t *testing.T) {

	t.Run("empty input channel", func(t *testing.T) {
		// Create an empty channel
		in := make(chan int)
		close(in)

		// Iterate over the batches and ensure no values are received
		for range BatchChanTimeout(in, 5, time.Hour) {
			t.Error("Unexpected batch received from an empty input channel")
		}
	})
	t.Run("max items per batch", func(t *testing.T) {
		in := make(chan int)
		go func() {
			defer close(in)
			for i := 0; i < 10; i++ {
				in <- i
			}
		}()

		var gotValues []int
		for ints := range BatchChanTimeout(in, 5, time.Hour) {
			require.Len(t, ints, 5)
			gotValues = append(gotValues, ints...)
		}

		require.Len(t, gotValues, 10)
		for i := 0; i < 10; i++ {
			require.Contains(t, gotValues, i)
		}
	})

	t.Run("max timeout", func(t *testing.T) {
		in := make(chan int)
		go func() {
			defer close(in)
			for i := 0; i < 10; i++ {
				in <- i
			}
		}()

		var gotValues []int
		for ints := range BatchChanTimeout(in, 100, time.Millisecond) {
			require.Len(t, ints, 10)
			gotValues = append(gotValues, ints...)
		}

		require.Len(t, gotValues, 10)
		for i := 0; i < 10; i++ {
			require.Contains(t, gotValues, i)
		}
	})
}

func TestManager(t *testing.T) {
	manager := NewManager[int](100, time.Second)
	go manager.ProcessMessages(nil)
	manager.Stop()
	err := manager.SendMessage(1)
	require.Error(t, err)
}
