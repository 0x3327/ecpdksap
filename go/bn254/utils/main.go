package utils

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)


func GenSECP256k1G1KeyPair() (privKey SECP256K1_fr.Element, pubKey SECP256K1.G1Affine) {

	privKey.SetRandom()

	var privKey_asBigInt big.Int
	privKey.BigInt(&privKey_asBigInt)

	pubKey.ScalarMultiplicationBase(&privKey_asBigInt)

	return privKey, pubKey
}

func GenG1KeyPair() (privKey fr.Element, pubKey bn254.G1Affine, _err error) {

	_, err := privKey.SetRandom()
	if err != nil {
		return fr.Element{}, bn254.G1Affine{}, fmt.Errorf("error generating private key: %w", err)
	}

	pubKeyAff, _ := CalcG1PubKey(privKey)

	return privKey, pubKeyAff, nil
}

func CalcG1PubKey(privKey fr.Element) (pubKey bn254.G1Affine, _err error) {

	privKeyBigInt := new(big.Int)
	privKey.BigInt(privKeyBigInt)

	g1Gen, _, _, _ := bn254.Generators()

	var pubKeyJac bn254.G1Jac
	pubKeyJac.ScalarMultiplication(&g1Gen, privKeyBigInt)

	var pubKeyAff bn254.G1Affine
	pubKeyAff.FromJacobian(&pubKeyJac)

	return pubKeyAff, nil
}

func GenG2KeyPair() (privKey fr.Element, pubKey bn254.G2Affine, _err error) {

	_, err := privKey.SetRandom()
	if err != nil {
		return fr.Element{}, bn254.G2Affine{}, fmt.Errorf("error generating private key: %w", err)
	}

	pubKeyAff, _ := CalcG2PubKey(privKey)

	return privKey, pubKeyAff, nil
}

func CalcG2PubKey(privKey fr.Element) (pubKey bn254.G2Affine, _err error) {

	privKeyBigInt := new(big.Int)
	privKey.BigInt(privKeyBigInt)

	_, g2Gen, _, _ := bn254.Generators()

	var pubKeyJac bn254.G2Jac
	pubKeyJac.ScalarMultiplication(&g2Gen, privKeyBigInt)

	var pubKeyAff bn254.G2Affine
	pubKeyAff.FromJacobian(&pubKeyJac)

	return pubKeyAff, nil
}

func Hash(input []byte) []byte {
	hasher := sha256.New()
	hasher.Write(input)     // Hash the input
	hash := hasher.Sum(nil) // Finalize the hash and return the result
	return hash
}

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

// computeStealthAddress computes the stealth address using pairings - from recipient perspective
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
	viewTagBytes := Hash(compressedBytesSlice)
	viewTag := viewTagBytes[0]

	return viewTag
}

func GenRandomRs(len int) (Rs []string, VTags []uint8) {

	for i := 0; i < len; i++ {
		r, R, _ := GenG1KeyPair()
		vTag := CalculateViewTag(&r, &R)
		Rs = append(Rs, R.X.String() + "." + R.Y.String())
		VTags = append(VTags, vTag)
	}

	return Rs, VTags
}


func UnpackXY (in string ) (X string, Y string) {
	separatorIdx:=	 strings.IndexByte(in, '.')
	X = in[:separatorIdx]
	Y = in[separatorIdx+1:]

	return
}