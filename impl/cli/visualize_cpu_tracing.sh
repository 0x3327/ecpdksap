# go test -v -bench=Benchmark_BN254_V2_V0_1byte_5000_Combined -run=1X -cpuprofile cpu.prof
go test -v -bench=Benchmark_BN254_5000 -run=1X -cpuprofile cpu.prof


# go tool pprof -http :8000 benchmark.test cpu.prof