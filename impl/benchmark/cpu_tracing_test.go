package benchmark

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
)

func Benchmark_BN254_V2_V0_1byte_1000(b *testing.B) {
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 1000)
	_Benchmark_BN254_V2_V0_1byte_ExpandedGnarkCrypto(b, 1000)
}

func Benchmark_BN254_V2_V0_1byte_5000_Combined(b *testing.B) {
	start := time.Now()
	Benchmark_BN254_V2_V0_1byte_Combined(b, 5000)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Bench total: ", elapsed)
}

func Benchmark_BN254_V2_V0_1byte_50000(b *testing.B) {
	start := time.Now()
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 50_000)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Bench total: ", elapsed)
}


func Benchmark_BN254_V2_V0_1byte_500000(b *testing.B) {
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 500_000)
}

func Benchmark_BN254_V2_V0_1byte_5000000(b *testing.B) {
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 5_000_000)
}

func _Benchmark_BN254_V2_V0_1byte_ExternalCalls(b *testing.B, sampleSize int) {

	hasher := sha256.New()
	combinedMeta, v_asBigIntPtr, _, K_SECP256k1_JacPtr := _generateData(sampleSize)
	var vR EC.G1Jac
	var vR_asAff EC.G1Affine
	_, _, _, g2Aff := EC.Generators()
	g2Aff_asArray := []EC.G2Affine{g2Aff}
	var Pv2_asJac SECP256K1.G1Jac
	
	b.ResetTimer()
	for _, cm := range combinedMeta {

		hasher.Reset()

		compressed := (vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))).Bytes()

		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(new (big.Int)))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExternalCalls) :: Total time:", b.Elapsed())
}

func _Benchmark_BN254_V2_V0_1byte_ExpandedGnarkCrypto(b *testing.B, sampleSize int) {

	hasher := sha256.New()
	combinedMeta, v_asBigIntPtr, _, K_SECP256k1_JacPtr := _generateData(sampleSize)
	var vR EC.G1Jac
	var vR_asAff EC.G1Affine
	_, _, _, g2Aff := EC.Generators()
	g2Aff_asArray := []EC.G2Affine{g2Aff}
	var Pv2_asJac SECP256K1.G1Jac

	neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix := EC.PrecomputationForFixedScalarMultiplication(v_asBigIntPtr)
	var table [15]EC.G1Jac

	b.ResetTimer()
	for _, cm := range combinedMeta {

		hasher.Reset()

		vR.FixedScalarMultiplication(cm.Rj, &table, neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix)

		compressed := (vR_asAff.FromJacobian(&vR)).Bytes()

		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(new (big.Int)))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExpandedGnarkCrypto) :: Total time:", b.Elapsed())
}