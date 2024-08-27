package benchmark

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	bls12_377 "ecpdksap-go/benchmark/curves/bls12-377"
	bls12_381 "ecpdksap-go/benchmark/curves/bls12-381"
	bls24_315 "ecpdksap-go/benchmark/curves/bls24-315"
	bn254 "ecpdksap-go/benchmark/curves/bn254"
	bw6_633 "ecpdksap-go/benchmark/curves/bw6-633"
	bw6_761 "ecpdksap-go/benchmark/curves/bw6-761"

	bn254_optimized "ecpdksap-go/benchmark/bn254"

	"ecpdksap-go/utils"
)

func RunBench(kind string, rndSeed int) {

	b := new(testing.B)
	b.StartTimer()

	if kind == "only-bn254" {

		bn254_optimized.Run(b, 5_000, 10, rndSeed)
		bn254_optimized.Run(b, 10_000, 10, rndSeed)
		bn254_optimized.Run(b, 20_000, 10, rndSeed)
		bn254_optimized.Run(b, 40_000, 10, rndSeed)
		bn254_optimized.Run(b, 80_000, 10, rndSeed)
		bn254_optimized.Run(b, 100_000, 10, rndSeed)

	} else if kind == "all-curves" {

		_Benchmark_Curves(b, 5_000, 10)
		// _Benchmark_Curves(b, 10_000, 10)
		// _Benchmark_Curves(b, 20_000, 10)
		// _Benchmark_Curves(b, 40_000, 10)
		// _Benchmark_Curves(b, 80_000, 10)
		// _Benchmark_Curves(b, 100_000, 10)
	}
}

func _Benchmark_Curves(b *testing.B, sampleSize int, nRepetitions int) {

	bls12_377.Run(b, sampleSize, nRepetitions, true)
	bls12_381.Run(b, sampleSize, nRepetitions, true)
	bls24_315.Run(b, sampleSize, nRepetitions, true)
	bn254.Run(b, sampleSize, nRepetitions, true)
	bw6_633.Run(b, sampleSize, nRepetitions, true)
	bw6_761.Run(b, sampleSize, nRepetitions, true)
}

func Benchmark_BN254_V2_V0_1byte_Combined(b *testing.B, sampleSize int) {

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

		compressed := (vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))).X.Bytes()

		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(new(big.Int)))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExternalCalls) :: Total time:", b.Elapsed())

	neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix := EC.PrecomputationForFixedScalarMultiplication(v_asBigIntPtr)
	var table [15]EC.G1Jac
	var a_El, b_El *EC_fp.Element
	var b_asBigInt big.Int

	b.ResetTimer()
	for _, cm := range combinedMeta {

		hasher.Reset()

		vR.FixedScalarMultiplication(cm.Rj, &table, neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix)

		a_El, b_El = vR_asAff.FromJacobianCoordX(&vR)

		compressed := (vR_asAff).X.Bytes()

		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			vR_asAff.FromJacobianCoordY(a_El, b_El, &vR)

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(&b_asBigInt))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExpandedGnarkCrypto) :: Total time:", b.Elapsed())
}

func _generateData(sampleSize int) (combinedMeta []*_CombinedMeta, v_asBigIntPtr *big.Int, V_Ptr EC.G1Jac, K_SECP256k1_JacPtr *SECP256K1.G1Jac) {
	_, v_asBigInt, V, _ := _EC_GenerateG1KeyPair()
	v_asBigIntPtr = &v_asBigInt
	V_Ptr = V

	//random data generation: Rj
	for j := 0; j < sampleSize; j++ {

		_, _, _, Rj_asAff := _EC_GenerateG1KeyPair()

		tmp := new(EC.G1Jac)
		tmp.FromAffine(&Rj_asAff)

		//note: store the last priv. key for R

		cm := new(_CombinedMeta)
		cm.Rj = new(EC.G1Jac)
		cm.Rj.FromAffine(&Rj_asAff)
		cm.Rj_asAffArr = []EC.G1Affine{Rj_asAff}
		cm.ViewTagTwoBytes = uint16(rand.Uint32() % 65536)
		cm.ViewTagSingleByte = uint8(rand.Uint32() % 256)

		combinedMeta = append(combinedMeta, cm)
	}

	_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
	var K_SECP256k1_Jac SECP256K1.G1Jac
	K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

	K_SECP256k1_JacPtr = &K_SECP256k1_Jac

	return
}

func _EC_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigInt big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigInt)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigInt)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

type _CombinedMeta struct {
	Rj                *EC.G1Jac
	Rj_asAffArr       []EC.G1Affine
	ViewTagTwoBytes   uint16
	ViewTagSingleByte uint8
}
