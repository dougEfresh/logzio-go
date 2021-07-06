package inMemoryQueue

import (
	queue "github.com/beeker1121/goque"
	"sync"
)

//Node storage of queue data
type Node struct {
	data []byte
	prev *Node
	next *Node
}

//ConcurrentQueue concurrent queue
type ConcurrentQueue struct {
	//mutex lock
	lock *sync.Mutex

	//empty and full locks
	notEmpty *sync.Cond
	notFull  *sync.Cond

	//queue storage backend
	backend *QueueBackend
}

func NewConcurrentQueue(maxSize uint64) *ConcurrentQueue {
	queue := ConcurrentQueue{}

	//init mutexes
	queue.lock = &sync.Mutex{}
	queue.notFull = sync.NewCond(queue.lock)
	queue.notEmpty = sync.NewCond(queue.lock)

	//init backend
	queue.backend = &QueueBackend{}
	queue.backend.size = 0
	queue.backend.head = nil
	queue.backend.tail = nil

	queue.backend.maxSize = maxSize
	return &queue
}

//QueueBackend Backend storage of the queue, a double linked list
type QueueBackend struct {
	//Pointers to root and end
	head *Node
	tail *Node

	//keep track of current size
	size    uint64
	maxSize uint64
}

func (queue *QueueBackend) createNode(data []byte) *Node {
	node := Node{}
	node.data = data
	node.next = nil
	node.prev = nil

	return &node
}

func (queue *QueueBackend) put(data []byte) error {

	if queue.size == 0 {
		//new root node
		node := queue.createNode(data)
		queue.head = node
		queue.tail = node

		queue.size += uint64(len(data))

		return nil
	}
	//queue non-empty append to head
	currentHead := queue.head
	newHead := queue.createNode(data)
	newHead.next = currentHead
	currentHead.prev = newHead

	queue.head = currentHead
	queue.size += uint64(len(data))
	return nil

}

func (queue *QueueBackend) pop() ([]byte, error) {

	currentEnd := queue.tail
	newEnd := currentEnd.prev

	if newEnd != nil {
		newEnd.next = nil
	}

	queue.size -= uint64(len(currentEnd.data))
	if queue.size == 0 {
		queue.head = nil
		queue.tail = nil
	}

	return currentEnd.data, nil
}

func (queue *QueueBackend) isEmpty() bool {
	return queue.size == 0
}

func (c *ConcurrentQueue) Enqueue(data []byte) (*Item, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	//insert
	err := c.backend.put(data)

	item := &Item{
		Value: data,
	}
	//signal notEmpty
	if err == nil {
		c.notEmpty.Signal()
	}
	return item, nil
}

type Item = queue.Item

func (c *ConcurrentQueue) Dequeue() (*Item, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for c.backend.isEmpty() {
		return nil, ErrEmpty
	}

	data, err := c.backend.pop()

	item := &Item{
		Value: data,
	}
	//signal notFull
	c.notFull.Signal()

	return item, err
}

func (c *ConcurrentQueue) Length() uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	// Size in bytes
	size := c.backend.size

	return size
}
func (c *ConcurrentQueue) Close() {

}
