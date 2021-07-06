package logzio

import queue "github.com/beeker1121/goque"

type Item = queue.Item

type genericQueue interface {
	Enqueue([]byte) (*Item, error)
	Dequeue() (*Item, error)
	Close()
	Length() uint64
}
