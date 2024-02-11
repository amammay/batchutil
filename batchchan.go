package batchutil

import (
	"fmt"
	"sync"
	"time"
)

// BatchChanTimeout will batch items from the input channel into slices of maxItems
// and send them to the returned channel. If maxTimeout is reached, the current
// batch will be sent and a new batch will be started.
func BatchChanTimeout[T any](in <-chan T, maxItems int, maxTimeout time.Duration) <-chan []T {
	batchC := make(chan []T)

	go func() {
		defer close(batchC)

		currentBatch := make([]T, 0, maxItems)
		timer := time.NewTimer(maxTimeout)
		defer timer.Stop()

		for {
			select {
			case value, ok := <-in:
				// channel is closed
				if !ok {
					if len(currentBatch) > 0 {
						batchC <- currentBatch
					}
					return
				}
				currentBatch = append(currentBatch, value)
				// check to see if we filled the max amount
				if len(currentBatch) == maxItems {
					timer.Stop()
					batchC <- currentBatch
					currentBatch = nil
					timer.Reset(maxTimeout)
				}
			case <-timer.C:
				if len(currentBatch) > 0 {
					batchC <- currentBatch
				}
				currentBatch = nil
				timer.Reset(maxTimeout)
			}
		}
	}()

	return batchC
}

// Manager is a batch manager that can be used to batch messages from many producers into a single consumer.
type Manager[T any] struct {
	// the channel to send messages to
	messageC chan T
	// wg is used to wait for all messages to be processed before shutting down
	senderWG sync.WaitGroup
	// processingDone is closed when the process is done
	processingDone    chan struct{}
	noMoreMessagesPlz chan struct{}

	batchSize     int
	rollingWindow time.Duration
}

// NewManager creates a new batch manager
// This is useful for batching messages from many producers into a single consumer.
// Experimental: This is a work in progress and may change.
func NewManager[T any](batchSize int, rollingWindow time.Duration) *Manager[T] {
	if rollingWindow == 0 {
		rollingWindow = 1 * time.Second
	}
	return &Manager[T]{messageC: make(chan T, 10000), processingDone: make(chan struct{}), noMoreMessagesPlz: make(chan struct{}), batchSize: batchSize, rollingWindow: rollingWindow}
}

// ProcessMessages starts the manager listening for messages and processing them.
// This is a blocking call and should be run in a goroutine.
func (m *Manager[T]) ProcessMessages(perBatch func([]T)) {
	for messages := range BatchChanTimeout(m.messageC, m.batchSize, m.rollingWindow) {
		if perBatch != nil {
			perBatch(messages)
		}
	}
	close(m.processingDone)
}

// Stop will stop the manager and wait for all messages to be processed.
func (m *Manager[T]) Stop() {
	// wait for all in-flight messages to be sent
	m.senderWG.Wait()

	// close the noMoreMessagesPlz channel to prevent any more messages from being sent
	close(m.noMoreMessagesPlz)
	// close the main message channel
	close(m.messageC)

	// wait for the process to be done
	<-m.processingDone
}

// SendMessage sends a message to the manager.
func (m *Manager[T]) SendMessage(message T) error {

	select {
	case <-m.noMoreMessagesPlz:
		return fmt.Errorf("no more messages please")
	default:
	}

	m.senderWG.Add(1)
	go func() {
		defer m.senderWG.Done()
		m.messageC <- message
	}()
	return nil
}
