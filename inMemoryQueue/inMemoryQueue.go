package inMemoryQueue

import (
	"fmt"
	queue "github.com/beeker1121/goque"
	"sync"
)

//ConcurrentQueue concurrent queue
type ConcurrentQueue struct {
	//mutex lock
	lock      *sync.RWMutex
	queue     chan []byte
	size      int
	maxLength int
	length    int
}

func NewConcurrentQueue(maxLength int) *ConcurrentQueue {
	queue := ConcurrentQueue{}
	//init mutexes
	queue.lock = &sync.RWMutex{}
	queue.queue = make(chan []byte, maxLength)
	queue.size = 0
	queue.length = 0
	queue.maxLength = maxLength
	return &queue
}

func (c *ConcurrentQueue) isEmpty() bool {
	c.lock.Lock()
	bool := c.size == 0
	c.lock.Unlock()
	return bool
}

func (c *ConcurrentQueue) Enqueue(data []byte) (*Item, error) {
	if !c.IsFull() {
		item := &Item{
			Value: data,
		}
		c.queue <- data
		c.lock.Lock()
		c.size += len(data)
		c.length++
		c.lock.Unlock()
		return item, nil
	}
	fmt.Printf("Queue is full dropping logs\n")
	return nil, nil
}

type Item = queue.Item

func (c *ConcurrentQueue) Dequeue() (*Item, error) {
	for c.isEmpty() {
		return nil, ErrEmpty
	}
	data := <-c.queue
	c.lock.Lock()
	c.size -= len(data)
	c.length--
	c.lock.Unlock()

	item := &Item{
		Value: data,
	}
	return item, nil
}

func (c *ConcurrentQueue) Length() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	size := c.size
	return uint64(size)
}
func (c *ConcurrentQueue) IsFull() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	isFull := c.maxLength-c.length == 0
	return isFull
}

func (c *ConcurrentQueue) Close() {
	var empty []byte
	for empty != nil {
		empty, _ := c.Dequeue()
		if empty == nil {
			break
		}
	}
	close(c.queue)
}
