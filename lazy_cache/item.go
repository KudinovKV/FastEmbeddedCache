package lazy_cache

import "time"

type item struct {
	value          any
	expirationDate time.Time
}

func (i *item) isExpired() bool {
	return i.expirationDate.Before(time.Now())
}

func newItem(value any, ttl time.Duration) *item {
	return &item{
		value:          value,
		expirationDate: time.Now().Add(ttl),
	}
}
