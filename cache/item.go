package cache

import "time"

type item struct {
	value       any
	expiredDate time.Time
}

func (i *item) isExpired() bool {
	return i.expiredDate.Before(time.Now())
}

func newItem(value any, ttl time.Duration) *item {
	return &item{
		value:       value,
		expiredDate: time.Now().Add(ttl),
	}
}
