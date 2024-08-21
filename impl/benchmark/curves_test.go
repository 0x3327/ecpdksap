package benchmark

import (
	bn254 "ecpdksap-go/benchmark/curves/bn254"
	"testing"
)

func Benchmark_BN254_5000(b *testing.B) {
	bn254.Run(b, 5_000, 10, true)
}

func Benchmark_BN254_1000000(b *testing.B) {
	bn254.Run(b, 1_000_000, 5, true)
}

func Benchmark_BN254_1000(b *testing.B) {
	bn254.Run(b, 1000, 1, true)
}

func Benchmark_Curves_10(b *testing.B) {
	_Benchmark_Curves(b, 10, 10)
}

func Benchmark_Curves_80000(b *testing.B) {
	_Benchmark_Curves(b, 80_000, 10)
}

