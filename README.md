# FastEmbeddedCache
Fast embedded key-value storage (redis-like) that supports Set, Get, and Delete methods with TTL.

## Install

```bash
go get -u github.com/KudinovKV/FastEmbededCache
```

## Example Usage: 

```go
import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/KudinovKV/FastEmbededCache/cache"
)


defaultTTL := time.Minute
cleanerTimeout := time.Minute*2

// create new instance with default ttl and cleaner timeout
driver := cache.New(context.Background(), defaultTTL, cleanerTimeout)

driver.Set("statue of liberty", "40.68960612218659, -74.0456618251789", time.Minute*2)

data, err := driver.Get("statue of liberty")
if err != nil {
    log.Fatal(err)
}

log.Println(data)

driver.Delete("statue of liberty")

_, err = driver.Get("statue of liberty")
if err != nil {
    if errors.Is(err, cache.ErrNotFound) {
        log.Println("element not found")
    }
}

driver.Shutdown()
```

## Features
- Used bare golang structures to guarantee `O(logN)` or `O(1)` complexity to retrieve requested elements.
- Expired data should be deleted automatically. TTL could be provided by a user, or default value could be used.
- Embedded solution (everything stored in RAM).
- Bare golang (didn't use external services, DBMS, or network data transfers within the algorithm).

## Algorithms

### Lazy Solution

Lazy solution uses bare Golang map and mutex to store data. This solution creates request to **DELETE** expired key only then this key is **GET** by client. Complexity for INSERT, GET, DELETE operation are `O(1)`. One of the problem here, that if there are a lot of key with small TTL, map would be full of expired entries and them will prevent other operations from executing more quickly.

### PriorityQueue Solution

This solution uses a binary tree from a `container/heap` to store "sorted" records. The main property of the binary heap is that it provides quick access to the element with the smallest TTL value without having to completely iterate over all the elements. In this algorithm, a separate goroutine is created, which goes through the heap and deletes expired keys. INSERT and DELETE operations become more complex. In the worst case, the complexity degrades to `O(log n)`. Additional memory is also required for heap storage. However, the problem of storing expired keys is solved, which can be critical on large data sets.

## Benchmark

Benchmark could be run using command below. Environment `TEST_MODE` set the algorithm to use: `LAZY` or `DEFAULT`.

```bash
make benchmark
```

### PriorityQueue Solution

| test name \ number of operation | 1000 (ns\op) | 10000 (ns\op) | 100000 (ns\op) | 1000000 (ns\op) | 10000000 (ns\op) |
|:-------------------------------:|:------------:|:-------------:|:--------------:|:---------------:|:----------------:|
|       BenchmarkReadSimple       |     111.8    |     111.7     |      175.8     |      209.8      |       236.9      |
|      BenchmarkReadParallel      |     249.8    |     316.5     |      242.6     |      211.7      |       246.2      |
|  BenchmarkReadParallelWithTTL   |     213.0    |     211.3     |      202.2     |      269.9      |       425.4      |
|      BenchmarkWriteSimple       |     233.2    |     452.9     |      395.2     |      381.1      |       460.0      |
|     BenchmarkWriteParallel      |     939.2    |      1123     |      1086      |       1206      |       5136       |

### Lazy Solution

| test name \ number of operation | 1000 (ns\op) | 10000 (ns\op) | 100000 (ns\op) | 1000000 (ns\op) | 10000000 (ns\op) |
|:-------------------------------:|:------------:|:-------------:|:--------------:|:---------------:|:----------------:|
|       BenchmarkReadSimple       |     110.5    |     137.1     |      141.1     |      214.5      |       235.8      |
|      BenchmarkReadParallel      |     274.7    |     320.5     |      215.9     |      213.6      |       216.9      |
|  BenchmarkReadParallelWithTTL   |     200.9    |     538.8     |      448.7     |      416.0      |       496.3      |
|      BenchmarkWriteSimple       |     259.4    |     213.9     |      291.2     |      374.2      |       429.3      |
|     BenchmarkWriteParallel      |     716.8    |     915.8     |      1076      |       1158      |       3467       |