package lazy_cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const defaultTTL = time.Minute

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

	driver := New(defaultTTL)

	for _, testCase := range testCases {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	for _, testCase := range testCases {
		actual, err := driver.Get(testCase.key)
		a.NoError(err)
		a.Equal(testCase.value, actual)
	}

	_, err := driver.Get("wrong key")
	a.Error(err)
	a.Equal(ErrNotFound, err)
}

func TestQueue_Delete(t *testing.T) {
	a := assert.New(t)

	driver := New(defaultTTL)

	for _, testCase := range testCases {
		driver.Set(testCase.key, testCase.value, testCase.ttl)
	}

	for _, testCase := range testCases {
		actual, err := driver.Get(testCase.key)
		a.NoError(err)
		a.Equal(testCase.value, actual)
	}

	value, err := driver.Get("3")
	a.NoError(err)
	a.Equal(3, value)

	value, err = driver.Get("4")
	a.NoError(err)
	a.Equal(4, value)

	driver.Delete("3")
	driver.Delete("4")
	driver.Delete("wrong key")

	_, err = driver.Get("3")
	a.Error(err)
	a.Equal(ErrNotFound, err)

	_, err = driver.Get("4")
	a.Error(err)
	a.Equal(ErrNotFound, err)
}
