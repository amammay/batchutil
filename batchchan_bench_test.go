package batchutil_test

import (
	"testing"
	"time"

	"github.com/amammay/batchutil"
)

func BenchmarkBatchChanTimeout(b *testing.B) {

	b.Run("unbuffered", func(b *testing.B) {

		// create an input channel with some test data
		in := make(chan int)
		go func() {
			defer close(in)
			for i := 0; i < b.N; i++ {
				in <- i
			}
		}()

		b.ResetTimer()
		for data := range batchutil.BatchChanTimeout(in, 100, time.Second) {
			_ = data
		}
	})

	b.Run("buffered", func(b *testing.B) {

		// create an input channel with some test data
		in := make(chan int, 1000)
		go func() {
			defer close(in)
			for i := 0; i < b.N; i++ {
				in <- i
			}
		}()

		b.ResetTimer()
		for data := range batchutil.BatchChanTimeout(in, 100, time.Second) {
			_ = data
		}
	})

}

func BenchmarkManager(b *testing.B) {

	b.Run("batch buffered", func(b *testing.B) {
		// create an input channel with some test data
		manager := batchutil.NewManager[int64](100, time.Second*5)

		go func() {
			for i := 0; i < b.N; i++ {
				manager.SendMessage(int64(i))
			}
			manager.Stop()
		}()

		b.ResetTimer()
		manager.ProcessMessages(func(_ []int64) {
			// do nothing
		})

	})
}
