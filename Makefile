COUNT=3
BENCH_TIME=1000x

.PHONY: read_benchmark
read_benchmark:
	go test -bench=BenchmarkRead. -benchmem -count=${COUNT} -benchtime ${BENCH_TIME}

.PHONY: write_benchmark
write_benchmark:
	go test -bench=BenchmarkWrite. -benchmem -count=${COUNT} -benchtime ${BENCH_TIME}

.PHONY: benchmark
benchmark: read_benchmark write_benchmark
