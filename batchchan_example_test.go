package batchutil

import (
	"fmt"
	"time"
)

func ExampleBatchChanTimeout() {

	// create an input channel with some test data
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
		}
	}()

	for data := range BatchChanTimeout(in, 2, time.Second*5) {
		fmt.Printf("%v\n", data)
	}

	// Output:
	// [0 1]
	// [2 3]
	// [4 5]
	// [6 7]
	// [8 9]

}

func ExampleBatchChanTimeout_window() {

	// create an input channel with some test data
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
		}
	}()

	for data := range BatchChanTimeout(in, 100, time.Second) {
		fmt.Printf("%v\n", data)
	}

	// all the output will be in one batch since it hit the one second window
	// Output:
	// [0 1 2 3 4 5 6 7 8 9]
}

func ExampleManager() {
	manager := NewManager[int](100, time.Second)
	for i := 0; i < 10; i++ {
		err := manager.SendMessage(i)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	go manager.ProcessMessages(func(vals []int) {

		for _, val := range vals {
			fmt.Printf("%v\n", val)
		}
	})

	manager.Stop()

	// Unordered output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}
