package inMemoryQueue

import (
	"fmt"
	"sync"
	"testing"
)

const (
	defaultQueueCapacity = 1000000 // 3 mb
)

func Test_CalculateSize(t *testing.T) {
	queue := NewConcurrentQueue(defaultQueueCapacity)
	message := make([]byte, 1000)
	queue.Enqueue(message)
	message = make([]byte, 1000)
	queue.Enqueue(message)
	if queue.Length() != 2000 && queue.length != 2 {
		t.Fatalf("Sum is inccorect.\nExpected: %d\nActual: %d", 2000, queue.Length())
	}
	queue.Dequeue()
	if queue.Length() != 1000 && queue.length != 1 {
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

func Test_FullQueue(t *testing.T) {
	queue := NewConcurrentQueue(1)
	message := make([]byte, 1000)
	queue.Enqueue(message)
	queue.Enqueue(message)
	item, _ := queue.Dequeue()
	item, _ = queue.Dequeue()
	if item != nil {
		t.Fatalf("Expected: none\nActual: %x", item)
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

func Test_Fifo(t *testing.T) {
	queue := NewConcurrentQueue(defaultQueueCapacity)

	queue.Enqueue(make([]byte, 12))
	queue.Enqueue(make([]byte, 10))
	queue.Enqueue(make([]byte, 10))
	queue.Enqueue(make([]byte, 10))
	queue.Enqueue(make([]byte, 10))
	_, err := queue.Dequeue()
	if err != nil {
		t.Fatalf("err")
	}
	if queue.Length() != 40 {
		t.Fatalf("Not fifo")
	}

}

func insertSamples(queue *ConcurrentQueue, wg *sync.WaitGroup) {
	for i := 0; i < 100000; i++ {
		_, err := queue.Enqueue(make([]byte, 1000))
		if err != nil {
		}
	}
	wg.Done()
}
func pullSamples(queue *ConcurrentQueue, wg *sync.WaitGroup) {
	for i := 0; i < 100000; i++ {
		_, err := queue.Dequeue()
		if err != nil {
		}
	}
	wg.Done()
}

func Test_EnqueueDequeueParallel(t *testing.T) {
	var wg sync.WaitGroup
	queue := NewConcurrentQueue(defaultQueueCapacity)
	for i := 0; i < 100000; i++ {
		_, err := queue.Enqueue(make([]byte, 10))
		if err != nil {
			t.Fatalf("en")
		}
	}
	wg.Add(1)
	go insertSamples(queue, &wg)
	wg.Add(1)
	go pullSamples(queue, &wg)
	wg.Wait()
	size := queue.Length()
	fmt.Printf("size: %d\n", size)
}

func BenchmarkEnqueueDequeueParallel(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		queue := NewConcurrentQueue(defaultQueueCapacity)
		for i := 0; i < 100000; i++ {
			_, err := queue.Enqueue(make([]byte, 1000))
			if err != nil {
				b.Fatalf("en")
			}
		}
		wg.Add(1)
		go insertSamples(queue, &wg)
		wg.Add(1)
		go pullSamples(queue, &wg)
		wg.Wait()
		size := queue.Length()
		if size != 100000000 {
			b.Fatalf("not the right size , %d", size)
		}
	}
}
func BenchmarkEnqueue500B(b *testing.B) {
	b.ReportAllocs()
	queue := NewConcurrentQueue(8000000)
	for i := 0; i < b.N; i++ {
		_, err := queue.Enqueue(make([]byte, 500))
		if err != nil {
			b.Fatalf("en")
		}

	}
	queue.Close()
}
func BenchmarkEnqueue10000B(b *testing.B) {
	b.ReportAllocs()
	queue := NewConcurrentQueue(defaultQueueCapacity)
	for i := 0; i < b.N; i++ {
		queue.Enqueue(make([]byte, 10000))
		queue.Dequeue()
	}
	queue.Close()
}

func BenchmarkDequeue500B(b *testing.B) {
	b.ReportAllocs()
	queue := NewConcurrentQueue(defaultQueueCapacity)
	for i := 0; i < b.N; i++ {
		_, err := queue.Enqueue(make([]byte, 500))
		_, err = queue.Dequeue()
		if err != nil {
			b.Fatalf("en")
		}
	}
}

func BenchmarkInitQueue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		queue := NewConcurrentQueue(1000000)
		queue.Close()
	}
}
