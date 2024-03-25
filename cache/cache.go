package cache

import (
	"container/heap"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/KudinovKV/FastEmbededCache/queue"
)

var ErrNotFound = errors.New("element not found")

// Driver implements default key-value interface with TTL logic.
// Driver runs cleaner goroutine on constructor and use context cancel function to stop it.
// To store input data Driver uses PriorityQueue based on heap.Interface.
type Driver struct {
	items map[string]*queue.Item
	queue queue.PriorityQueue
	mu    sync.RWMutex

	cleanerTimeout time.Duration
	defaultTTL     time.Duration

	cancel context.CancelFunc
	stopCh chan struct{}
}

func (d *Driver) Set(key string, value any, ttl time.Duration) {
	if ttl == 0 {
		ttl = d.defaultTTL
	}

	expirationDate := time.Now().Add(ttl)
	newItem := &queue.Item{
		Key:            key,
		Value:          value,
		ExpirationDate: expirationDate,
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if oldItem, ok := d.items[key]; ok {
		oldItem.ExpirationDate = expirationDate
		oldItem.Value = value

		heap.Fix(&d.queue, oldItem.Index)

		return
	}

	heap.Push(&d.queue, newItem)
	d.items[key] = newItem
}

func (d *Driver) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if item, ok := d.items[key]; ok {
		heap.Remove(&d.queue, item.Index)
		delete(d.items, key)
	}
}

func (d *Driver) Get(key string) (any, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	item, ok := d.items[key]
	if !ok {
		return nil, ErrNotFound
	}

	if item.IsExpired() {
		return nil, ErrNotFound
	}

	return item.Value, nil
}

func (d *Driver) Shutdown() {
	d.cancel()

	<-d.stopCh
}

func (d *Driver) runCleaner(ctx context.Context) {
	ticker := time.NewTicker(d.cleanerTimeout)

	for {
		select {
		case <-ticker.C:
			d.deleteExpiredKeys()
		case <-ctx.Done():
			ticker.Stop()

			d.stopCh <- struct{}{}

			return
		}
	}
}

// deleteExpiredKeys iterate through sorted d.queue. Starting from the element which expired earlier.
// If Head element is not expired, we don't need to check another elements.
func (d *Driver) deleteExpiredKeys() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for d.queue.Len() > 0 && d.queue.Head().(*queue.Item).IsExpired() {
		expiredItem := heap.Pop(&d.queue).(*queue.Item)

		delete(d.items, expiredItem.Key)
	}
}

func New(ctx context.Context, defaultTTL, cleanerTimeout time.Duration) *Driver {
	d := &Driver{
		items: make(map[string]*queue.Item),
		queue: make(queue.PriorityQueue, 0),

		cleanerTimeout: cleanerTimeout,
		defaultTTL:     defaultTTL,

		stopCh: make(chan struct{}),
	}

	heap.Init(&d.queue)

	ctx, d.cancel = context.WithCancel(ctx)
	go d.runCleaner(ctx)

	return d
}
