package v2

import (
	"ecpdksap-go/utils"
	"encoding/hex"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"

	// "github.com/ethereum/go-ethereum/crypto/sha3"
	"golang.org/x/crypto/sha3"
)

// Computes the shared secret - from sender's perspective
func SenderComputesSharedSecret(r *fr.Element, V *bn254.G1Affine, K *SECP256K1.G1Affine) (bn254.GT) {

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

	var res bn254.E2

	res.Add(&pubKey.C0.B0, &pubKey.C0.B1)
	res.Add(&res, &pubKey.C1.B0)
	res.Add(&res, &pubKey.C1.B1)

	res.A0.BigInt(&b)

	b.Add(&b, res.A1.BigInt(new(big.Int)))

	return
}


func Compute_b_asElement(pubKey *bn254.GT) (b SECP256K1_fr.Element) {

	var res bn254.E2

	res.Add(&pubKey.C0.B0, &pubKey.C0.B1)
	res.Add(&res, &pubKey.C1.B0)
	res.Add(&res, &pubKey.C1.B1)

	b = SECP256K1_fr.Element(res.A0)
	b2 :=SECP256K1_fr.Element(res.A1)
	
	b.Add(&b,  &b2)

	return
}

func SenderComputesEthAddress(b *SECP256K1_fr.Element, K *SECP256K1.G1Affine) (string) {

	var b_asBigInt big.Int
	b.BigInt(&b_asBigInt)

	var P SECP256K1.G1Affine
	P.ScalarMultiplication(K, &b_asBigInt)

	return ComputeEthAddress(&P)
}


// Computes the shared secred - from recipents's perspective
func RecipientComputesSharedSecret(v *fr.Element, R *bn254.G1Affine, K2 *SECP256K1.G1Affine) (bn254.GT) {

	productAffine := utils.BN254_MulG1PointandElement(*R, *v)

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

func ComputeEthAddress(P *SECP256K1.G1Affine) (addr string){

	Px_Bytes := P.X.Bytes()
	Py_Bytes := P.Y.Bytes()
	hash := sha3.NewLegacyKeccak256()

	hash.Write(Px_Bytes[:])
	hash.Write(Py_Bytes[:])
	buf := hash.Sum(nil)

	totalLen := len(buf)
	addr = "0x"+ hex.EncodeToString(buf[totalLen-20:])
	return
}