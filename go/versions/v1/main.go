package v1

import (
	"fmt"
	"math/big"

	BN254 "github.com/consensys/gnark-crypto/ecc/bn254"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"ecpdksap-go/utils"
)

// computeStealthAddress computes the stealth address using pairings - from sender perspective
func SenderComputesStealthPubKey(r *BN254_fr.Element, V *BN254.G1Affine, K *BN254.G2Affine) (BN254.GT, error) {
	rBigInt := new(big.Int)
	r.BigInt(rBigInt)

	var V_Jac BN254.G1Jac
	var rV_product BN254.G1Jac
	rV_product.ScalarMultiplication(V_Jac.FromAffine(V), rBigInt)

	var productAffine BN254.G1Affine
	productAffine.FromJacobian(&rV_product)

	hash_asBytes := utils.BN254_HashG1Point(&productAffine)
	var hash BN254_fr.Element
	hash.SetBytes(hash_asBytes)
	var hash_asBigInt big.Int
	hash.BigInt(&hash_asBigInt)

	var g1Point BN254.G1Affine
	g1Point.ScalarMultiplicationBase(&hash_asBigInt)

	P, err := BN254.Pair([]BN254.G1Affine{g1Point}, []BN254.G2Affine{*K})
	if err != nil {
		return BN254.GT{}, fmt.Errorf("error computing pairing: %w", err)
	}

	return P, nil
}

// computes the stealth public key using pairings - from recipient perspective
func RecipientComputesStealthPubKey(k *BN254_fr.Element, v *BN254_fr.Element, R *BN254.G1Affine) (BN254.GT) {

	vBigInt := new(big.Int)
	v.BigInt(vBigInt)

	var vR_product BN254.G1Affine
	vR_product.ScalarMultiplication(R, vBigInt)

	hash_asBytes := utils.BN254_HashG1Point(&vR_product)
	var hash BN254_fr.Element
	hash.SetBytes(hash_asBytes)
	
	var privKey BN254_fr.Element
	privKey.Mul(&hash, k)
	var privKey_asBigInt big.Int
	privKey.BigInt(&privKey_asBigInt)

	_, _, g1Aff, g2Aff := BN254.Generators()
	
	pairingResult, _ := BN254.Pair([]BN254.G1Affine{g1Aff}, []BN254.G2Affine{g2Aff})

	var P BN254.GT
	P.CyclotomicExp(pairingResult, &privKey_asBigInt)

	return P
}

func ViewerComputesStealthPubKey(K *BN254.G2Affine, R *BN254.G1Affine, v *BN254_fr.Element) (BN254.GT) {

	vBigInt := new(big.Int)
	v.BigInt(vBigInt)

	var rV_product BN254.G1Affine
	rV_product.ScalarMultiplication(R, vBigInt)

	hash_asBytes := utils.BN254_HashG1Point(&rV_product)
	var hash BN254_fr.Element
	hash.SetBytes(hash_asBytes)
	var hash_asBigInt big.Int
	hash.BigInt(&hash_asBigInt)

	var g1Point BN254.G1Affine
	g1Point.ScalarMultiplicationBase(&hash_asBigInt)
	
	pairingResult, _ := BN254.Pair([]BN254.G1Affine{g1Point}, []BN254.G2Affine{*K})

	return pairingResult
}


func CalculateViewTag(r *BN254_fr.Element, V *BN254.G1Affine) uint8 {
	// Convert r to big.Int
	rBigInt := new(big.Int)
	r.BigInt(rBigInt)

	// Perform scalar multiplication of V by r
	var VJac BN254.G1Jac
	var product BN254.G1Jac
	product.ScalarMultiplication(VJac.FromAffine(V), rBigInt)

	// Convert the product to compressed bytes
	var productAffine BN254.G1Affine
	compressedBytes := productAffine.FromJacobian(&product).Bytes()

	// Convert [64]byte array to slice
	compressedBytesSlice := compressedBytes[:]

	// Extract the first byte of the hashed field element as the view tag
	viewTagBytes := utils.Hash(compressedBytesSlice)
	viewTag := viewTagBytes[0]

	return viewTag
}