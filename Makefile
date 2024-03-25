COUNT=1
BENCH_TIME=100000x

# for lazy algorithm
# TEST_MODE = "LAZY"
TEST_MODE ?= "LAZY"

.PHONY: read_benchmark
read_benchmark:
	TEST_MODE=${TEST_MODE} go test -bench=BenchmarkRead. -benchmem -count=${COUNT} -benchtime=${BENCH_TIME}

.PHONY: write_benchmark
write_benchmark:
	TEST_MODE=${TEST_MODE} go test -bench=BenchmarkWrite. -benchmem -count=${COUNT} -benchtime=${BENCH_TIME}

.PHONY: benchmark
benchmark: read_benchmark write_benchmark

.PHONY: test
test:
	go test -race -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...