package batchutil

import (
	"errors"
)

// BatchFunc is called for each batch.
// Any error will cancel the batching operation but returning ErrAbort
// indicates it was deliberate, and not an error case.
type BatchFunc[T any] func([]T) error

// ErrAbort indicates that the operation was aborted deliberately.
var ErrAbort = errors.New("done")

// All calls eachFn for all items
// Returns any error from eachFn except for Abort it returns nil.
func All[T any](data []T, batchSize int, eachFn BatchFunc[T]) error {
	count := len(data)
	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}
		batch := make([]T, end-i)
		copy(batch, data[i:end])
		err := eachFn(batch)
		if errors.Is(err, ErrAbort) {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}
