package benchmark

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	"ecpdksap-go/utils"
)

func Benchmark_BN254_V2_V0_1byte_1000(b *testing.B) {
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 1000)
}

func Benchmark_BN254_V2_V0_1byte_5000(b *testing.B) {
	start := time.Now()
	_Benchmark_BN254_V2_V0_1byte_ExternalCalls(b, 5000)
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

func _generateData (sampleSize int) (combinedMeta  []*_CombinedMeta, v_asBigIntPtr *big.Int, V_Ptr EC.G1Jac, K_SECP256k1_JacPtr *SECP256K1.G1Jac) {
	_, v_asBigInt, V, _ := _EC_GenerateG1KeyPair()
	v_asBigIntPtr = &v_asBigInt
	V_Ptr = V

	var r_asBigInt big.Int

	//random data generation: Rj
	var Rs []EC.G1Jac
	var Rs_Ptr []*EC.G1Jac
	var RsAff_asArr [][]EC.G1Affine

	var rs []big.Int

	for j := 0; j < sampleSize; j++ {

		_, rj_asBigInt, Rj, Rj_asAff := _EC_GenerateG1KeyPair()

		Rs = append(Rs, Rj)
		RsAff_asArr = append(RsAff_asArr, []EC.G1Affine{Rj_asAff})

		tmp := new(EC.G1Jac)
		tmp.FromAffine(&Rj_asAff)
		Rs_Ptr = append(Rs_Ptr, tmp)

		//note: store the last priv. key for R
		r_asBigInt = rj_asBigInt
		rs = append(rs, r_asBigInt)

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


func _EC_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}


type _CombinedMeta struct {
	Rj                *EC.G1Jac
	Rj_asAffArr       []EC.G1Affine
	ViewTagTwoBytes   uint16
	ViewTagSingleByte uint8
}
