go test -v -bench=Benchmark_BN254_V2_V0_1byte_50000 -run=1X -cpuprofile cpu.prof

go tool pprof -http :8000 benchmark.test cpu.prof