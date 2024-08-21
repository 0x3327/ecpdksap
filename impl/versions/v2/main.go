package v2

import (
	"encoding/hex"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"

	// "github.com/ethereum/go-ethereum/crypto/sha3"
	"golang.org/x/crypto/sha3"
)

// Computes the shared secret - from sender's perspective
func SenderComputesSharedSecret(r *fr.Element, V *bn254.G1Affine, K *SECP256K1.G1Affine) bn254.GT {

	rBigInt := new(big.Int)
	r.BigInt(rBigInt)

	// Perform scalar multiplication of V by r
	var V_Jac bn254.G1Jac
	var rV_product bn254.G1Jac
	rV_product.ScalarMultiplication(V_Jac.FromAffine(V), rBigInt)

	var productAffine bn254.G1Affine
	productAffine.FromJacobian(&rV_product)

	// Compute pairing
	_, g2Gen, _, _ := bn254.Generators()

	one := new(big.Int)
	one.SetString("1", 10)

	var G2Jac bn254.G2Jac
	G2Jac.ScalarMultiplication(&g2Gen, one)

	var G2Aff bn254.G2Affine
	G2Aff.FromJacobian(&G2Jac)

	P, _ := bn254.Pair([]bn254.G1Affine{productAffine}, []bn254.G2Affine{G2Aff})

	return P
}

func Compute_b(pubKey *bn254.GT) (b big.Int) {

	return *pubKey.C0.B0.A0.BigInt(new(big.Int))
}

func Compute_b_asElement(pubKey *bn254.GT) (b SECP256K1_fr.Element) {

	b_asBigInt := Compute_b(pubKey)

	return *b.SetBigInt(&b_asBigInt)
}

func SenderComputesEthAddress(b *SECP256K1_fr.Element, K *SECP256K1.G1Affine) string {

	var b_asBigInt big.Int
	b.BigInt(&b_asBigInt)

	var P SECP256K1.G1Affine
	P.ScalarMultiplication(K, &b_asBigInt)

	return ComputeEthAddress(&P)
}

// Computes the shared secred - from recipents's perspective
func RecipientComputesSharedSecret(v *fr.Element, R *bn254.G1Affine, K2 *SECP256K1.G1Affine) bn254.GT {
	_, _, _, g2Aff := EC.Generators()

	neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix := EC.PrecomputationForFixedScalarMultiplication(v.BigInt(new(big.Int)))
	var table [15]EC.G1Jac

	precomputedQLines := [][2][66]EC.LineEvaluationAff {EC.PrecomputeLines(g2Aff)}

	var vR, R_asJac EC.G1Jac
	var vR_asAff EC.G1Affine
	vR.FixedScalarMultiplication(R_asJac.FromAffine(R), &table, neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix)
	
	P, _ := EC.PairFixedQ([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, precomputedQLines)

	return P
}

func ComputeEthAddress(P *SECP256K1.G1Affine) (addr string) {

	Px_Bytes := P.X.Bytes()
	Py_Bytes := P.Y.Bytes()
	hash := sha3.NewLegacyKeccak256()

	hash.Write(Px_Bytes[:])
	hash.Write(Py_Bytes[:])
	buf := hash.Sum(nil)

	totalLen := len(buf)
	addr = "0x" + hex.EncodeToString(buf[totalLen-20:])
	return
}
