package batchutil

import "fmt"

func ExampleAll() {

	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	err := All(data, 3, func(batch []string) error {
		fmt.Printf("%v\n", batch)
		// Do something with the batch
		return nil
	})
	if err != nil {
		// Handle the error
	}

	// Output:
	// [a b c]
	// [d e f]
	// [g h i]
	// [j]
}

func ExampleAll_abort() {

	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	err := All(data, 3, func(batch []string) error {
		fmt.Printf("%v\n", batch)
		// Do something with the batch
		return ErrAbort
	})
	if err != nil {
		// Handle the error
	}
	fmt.Println("aborted")

	// Output:
	// [a b c]
	// aborted
}
