package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/KudinovKV/FastEmbededCache/cache"
	"github.com/KudinovKV/FastEmbededCache/lazy_cache"
)

const (
	testModeLazy = "LAZY"

	defaultTTL     = time.Minute
	cleanerTimeout = time.Second
)

type cacheDriver interface {
	Set(key string, value any, ttl time.Duration)
	Get(key string) (any, error)
	Delete(key string)
}

func initDriver() cacheDriver {
	mode := os.Getenv("TEST_MODE")

	switch mode {
	case testModeLazy:
		return lazy_cache.New(defaultTTL)
	default:
		return cache.New(context.Background(), defaultTTL, cleanerTimeout)
	}
}

func BenchmarkWriteSimple(b *testing.B) {
	driver := initDriver()

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Duration(rand.Int()%100)*time.Second)
	}
}

func BenchmarkWriteParallel(b *testing.B) {
	driver := initDriver()

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(i int) {
			driver.Set(fmt.Sprintf("%d", i), i, time.Duration(rand.Int()%100)*time.Second)

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkReadSimple(b *testing.B) {
	driver := initDriver()

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Duration(rand.Int()%100)*time.Minute)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := driver.Get(fmt.Sprintf("%d", i))
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkReadParallel(b *testing.B) {
	driver := initDriver()

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Duration(rand.Int()%100)*time.Minute)
	}

	b.ResetTimer()

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(i int) {
			_, err := driver.Get(fmt.Sprintf("%d", i))
			if err != nil {
				b.Error(err)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkReadParallelWithTTL(b *testing.B) {
	driver := initDriver()

	for i := 0; i < b.N; i++ {
		ttl := time.Minute
		if i%5 == 0 {
			ttl = time.Millisecond
		}

		driver.Set(fmt.Sprintf("%d", i), i, ttl)
	}

	b.ResetTimer()

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(i int) {
			_, err := driver.Get(fmt.Sprintf("%d", i))
			if err != nil {
				if i%5 != 0 {
					b.Error(err)
				}
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
}
