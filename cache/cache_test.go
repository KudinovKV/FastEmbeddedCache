package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultTTL            = time.Minute
	defaultCleanerTimeout = time.Minute
)

var testCases = []struct {
	key   string
	value int
	ttl   time.Duration
}{
	{"1", 1, time.Second},
	{"4", 4, time.Minute},
	{"3", 3, time.Second * 30},
	{"2", 2, time.Second * 15},
	{"5", 5, time.Minute + time.Second},
}

func TestQueue_SetAndGet(t *testing.T) {
	a := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	driver := New(ctx, defaultTTL, defaultCleanerTimeout)

	for _, testCase := range testCases {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	item, err := driver.Get("1")
	a.NoError(err)
	a.Equal(1, item)

	item, err = driver.Get("2")
	a.NoError(err)
	a.Equal(2, item)

	item, err = driver.Get("3")
	a.NoError(err)
	a.Equal(3, item)

	item, err = driver.Get("4")
	a.NoError(err)
	a.Equal(4, item)

	item, err = driver.Get("5")
	a.NoError(err)
	a.Equal(5, item)

	_, err = driver.Get("wrong key")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	time.Sleep(time.Second + time.Millisecond)

	_, err = driver.Get("1")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	_, err = driver.Get("5")
	a.NoError(err)
	a.Equal(5, item)

	_, err = driver.Get("wrong key")
	a.Error(err)
	a.Equal(ErrNotFound, err)
}

func TestQueue_SameSet(t *testing.T) {
	a := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	driver := New(ctx, defaultTTL, time.Second)

	for _, testCase := range testCases {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	item, err := driver.Get("1")
	a.NoError(err)
	a.Equal(1, item)

	driver.Set("1", "new value", time.Minute)

	item, err = driver.Get("1")
	a.NoError(err)
	a.Equal("new value", item)

	time.Sleep(time.Second + time.Millisecond)

	item, err = driver.Get("1")
	a.NoError(err)
	a.Equal("new value", item)

}

func TestQueue_Delete(t *testing.T) {
	a := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	driver := New(ctx, defaultTTL, defaultCleanerTimeout)

	for _, testCase := range testCases {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	item, err := driver.Get("3")
	a.NoError(err)
	a.Equal(3, item)

	item, err = driver.Get("4")
	a.NoError(err)
	a.Equal(4, item)

	driver.Delete("3")
	driver.Delete("4")
	driver.Delete("wrong key")

	_, err = driver.Get("3")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	_, err = driver.Get("4")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	_, err = driver.Get("wrong key")
	a.Error(err)
	a.Equal(ErrNotFound, err)
}

func TestQueue_Cleaner(t *testing.T) {
	a := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	driver := New(ctx, defaultTTL, time.Second)

	testCasesForCleaner := []struct {
		key   string
		value int
		ttl   time.Duration
	}{
		{"1", 1, time.Millisecond},
		{"4", 4, time.Minute},
		{"3", 3, time.Second * 30},
		{"2", 2, time.Millisecond * 5},
		{"5", 5, time.Minute + time.Second},
	}

	for _, testCase := range testCasesForCleaner {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	item, err := driver.Get("1")
	a.NoError(err)
	a.Equal(1, item)

	item, err = driver.Get("2")
	a.NoError(err)
	a.Equal(2, item)

	item, err = driver.Get("5")
	a.NoError(err)
	a.Equal(5, item)

	time.Sleep(time.Second * 2)

	_, err = driver.Get("1")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	_, err = driver.Get("2")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	item, err = driver.Get("5")
	a.NoError(err)
	a.Equal(5, item)
}
