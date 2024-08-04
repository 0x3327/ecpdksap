package benchmark

import (
	"testing"

	bn254 "ecpdksap-go/benchmark/curves/bn254"
)

func main(b *testing.B) {
	bn254.Run(b, 5_000, 10, true)
}
