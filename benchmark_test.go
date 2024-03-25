package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/KudinovKV/FastEmbededCache/cache"
)

func BenchmarkWriteSimple(b *testing.B) {
	driver := cache.New(time.Minute)

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Duration(i)*time.Second)
	}
}

func BenchmarkWriteParallel(b *testing.B) {
	driver := cache.New(time.Minute)

	wg := sync.WaitGroup{}
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func(i int) {
			driver.Set(fmt.Sprintf("%d", i), i, time.Duration(i)*time.Second)

			wg.Done()
		}(i)
	}

	wg.Wait()
}

func BenchmarkReadSimple(b *testing.B) {
	driver := cache.New(time.Minute)

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Minute)
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
	driver := cache.New(time.Minute)

	for i := 0; i < b.N; i++ {
		driver.Set(fmt.Sprintf("%d", i), i, time.Minute)
	}

	b.ResetTimer()

	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
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
	driver := cache.New(time.Minute)

	for i := 0; i < b.N; i++ {
		ttl := time.Minute
		if i%5 == 0 {
			ttl = time.Millisecond
		}

		driver.Set(fmt.Sprintf("%d", i), i, ttl)
	}

	b.ResetTimer()

	var wg sync.WaitGroup

	for i := 0; i < b.N; i++ {
		wg.Add(1)
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
