package benchmark

import (
	"testing"

	bn254_bench "ecpdksap-go/benchmark/curves/bn254"
)

func Benchmark_Curves(b *testing.B) {

	bn254_bench.Run(b, 5_000, 3, true)

}
