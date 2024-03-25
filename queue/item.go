package queue

import "time"

// An Item is something we manage in a priority queue.
type Item struct {
	Key            string
	Value          any
	ExpirationDate time.Time // The priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
}

func (i *Item) IsExpired() bool {
	return i.ExpirationDate.Before(time.Now())
}
