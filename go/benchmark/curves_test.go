package benchmark

import (
	"testing"

	bls12_377 "ecpdksap-go/benchmark/curves/bls12-377"
	bls12_381 "ecpdksap-go/benchmark/curves/bls12-381"
	bls24_315 "ecpdksap-go/benchmark/curves/bls24-315"
	bn254 "ecpdksap-go/benchmark/curves/bn254"
	bw6_633 "ecpdksap-go/benchmark/curves/bw6-633"
	bw6_761 "ecpdksap-go/benchmark/curves/bw6-761"
)

func Benchmark_BN254_5000(b *testing.B) {
	bn254.Run(b, 5_000, 10, true)
}

func Benchmark_BN254_1000(b *testing.B) {
	bn254.Run(b, 1000, 1, true)
}

func Benchmark_Curves_10(b *testing.B) {
	_Benchmark_Curves(b, 10, 3)
}

func Benchmark_Curves_80000(b *testing.B) {
	_Benchmark_Curves(b, 80_000, 10)
}

func _Benchmark_Curves(b *testing.B, sampleSize int, nRepetitions int) {

	bls12_377.Run(b, sampleSize, nRepetitions, true)
	bls12_381.Run(b, sampleSize, nRepetitions, true)
	bls24_315.Run(b, sampleSize, nRepetitions, true)
	bn254.Run(b, sampleSize, nRepetitions, true)
	bw6_633.Run(b, sampleSize, nRepetitions, true)
	bw6_761.Run(b, sampleSize, nRepetitions, true)
}
