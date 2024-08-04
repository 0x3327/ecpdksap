package benchmark

import (
	"fmt"
	"math/big"
	"testing"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func Benchmark_ScalarMultiplication(b *testing.B) {
	_, _, V, _ := _EC_GenerateG1KeyPair()
	_, rnd_asBigInt, _, _ := _EC_GenerateG1KeyPair()
	var RND EC.G1Jac

	b.ResetTimer()

	for i := 0; i < 5000; i++ {
		RND.ScalarMultiplication(&V, &rnd_asBigInt)
	}

	fmt.Println("Elapsed:", b.Elapsed())
}

func Benchmark_BatchScalarMultiplication(b *testing.B) {
	_, _, V, _ := _EC_GenerateG1KeyPair()
	_, rnd_asBigInt, _, _ := _EC_GenerateG1KeyPair()
	var RND EC.G1Jac

	var Vs []*EC.G1Jac

	for i := 0; i < 5000; i++ {
		// RND.ScalarMultiplication(&V, &rnd_asBigInt)
		Vs = append(Vs, &V)
	}


	b.ResetTimer()

	RND.BatchScalarMultiplicationUsingFixedS(&Vs, &rnd_asBigInt)

	fmt.Println("Elapsed:", b.Elapsed())
}


func _EC_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}
