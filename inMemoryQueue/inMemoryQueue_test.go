package inMemoryQueue

import (
	"testing"
)

const (
	defaultQueueCapacity = 3 * 1024 * 1024 // 3 mb
)

func Test_CalculateSize(t *testing.T) {
	queue := NewConcurrentQueue(defaultQueueCapacity)
	message := make([]byte, 1000)
	queue.Enqueue(message)
	message = make([]byte, 1000)
	queue.Enqueue(message)
	if queue.Length() != 2000 {
		t.Fatalf("Sum is inccorect.\nExpected: %d\nActual: %d", 2000, queue.Length())
	}
	queue.Dequeue()
	if queue.Length() != 1000 {
		t.Fatalf("Sum is inccorect.\nExpected: %d\nActual: %d", 1000, queue.Length())
	}
}

func Test_EmptyQueue(t *testing.T) {
	queue := NewConcurrentQueue(defaultQueueCapacity)
	_, err := queue.Dequeue()
	if err == nil {
		t.Fatalf("Expected Error: %s\nActual: none", err)
	}
	message := make([]byte, 1000)
	queue.Enqueue(message)
	_, err = queue.Dequeue()
	if err != nil {
		t.Fatalf("Expected: none\nActual: %s", err)
	}
}

func Test_EnqueueDequeue(t *testing.T) {
	queue := NewConcurrentQueue(defaultQueueCapacity)
	item, err := queue.Enqueue([]byte("sample-message"))
	if err != nil {
		t.Fatal()
	}
	if item == nil {
		t.Fatalf("Item is nil")
	}
	item, err = queue.Dequeue()
	if err != nil {
		t.Fatal()
	}
	if item == nil {
		t.Fatalf("Item is nil")
	}

}
