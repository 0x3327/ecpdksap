package v2

import (
	"ecpdksap-bn254/utils"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// computeStealthAddress computes the stealth address using pairings - from sender perspective
func SenderComputesStealthPubKey(r *fr.Element, V *bn254.G1Affine, K *bn254.G2Affine) (bn254.GT, error) {
	// Convert vPrivateKey to big.Int for cyclotomic exponentiation
	rBigInt := new(big.Int)
	r.BigInt(rBigInt)

	// Perform scalar multiplication of V by r
	var V_Jac bn254.G1Jac
	var rV_product bn254.G1Jac
	rV_product.ScalarMultiplication(V_Jac.FromAffine(V), rBigInt)

	// Convert the product to compressed bytes
	var productAffine bn254.G1Affine
	productAffine.FromJacobian(&rV_product)

	// Compute pairing
	P, err := bn254.Pair([]bn254.G1Affine{productAffine}, []bn254.G2Affine{*K})
	if err != nil {
		return bn254.GT{}, fmt.Errorf("error computing pairing: %w", err)
	}

	return P, nil
}

// computes the stealth public key using pairings - from recipient perspective
func RecipientComputesStealthPubKey(K *bn254.G2Affine, R *bn254.G1Affine, v *fr.Element) (bn254.GT, error) {
	// Convert vPrivateKey to big.Int for cyclotomic exponentiation
	vBigInt := new(big.Int)
	v.BigInt(vBigInt)
	// Compute pairing
	pairingResult, err := bn254.Pair([]bn254.G1Affine{*R}, []bn254.G2Affine{*K})
	// fmt.Println("pairingResult in bytes:", pairingResult.Bytes())
	if err != nil {
		return bn254.GT{}, fmt.Errorf("error computing pairing: %w", err)
	}

	// Compute cyclotomic exponentiation
	var P bn254.GT
	P.CyclotomicExp(pairingResult, vBigInt)

	return P, nil
}

func CalculateViewTag(r *fr.Element, V *bn254.G1Affine) uint8 {
	// Convert r to big.Int
	rBigInt := new(big.Int)
	r.BigInt(rBigInt)

	// Perform scalar multiplication of V by r
	var VJac bn254.G1Jac
	var product bn254.G1Jac
	product.ScalarMultiplication(VJac.FromAffine(V), rBigInt)

	// Convert the product to compressed bytes
	var productAffine bn254.G1Affine
	compressedBytes := productAffine.FromJacobian(&product).Bytes()

	// Convert [64]byte array to slice
	compressedBytesSlice := compressedBytes[:]

	// Extract the first byte of the hashed field element as the view tag
	viewTagBytes := utils.Hash(compressedBytesSlice)
	viewTag := viewTagBytes[0]

	return viewTag
}