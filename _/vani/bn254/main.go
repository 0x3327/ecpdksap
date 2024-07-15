package main

import (
	"fmt"
	"math/big"
	"sap-go/config"
	"time"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func testSearchSpeed() {
	kPrivateKey, _ := generatePrivateKey()
	vPrivateKey, _ := generatePrivateKey()
	rPrivateKey, _ := generatePrivateKey()

	kPublicKey, _, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)

	_, g2Gen, _, _ := bn254.Generators()
	publicKeys := make([]bn254.G2Affine, 0, config.RunNumber+1)
	for i := 0; i < config.RunNumber; i++ {
		randomPrivateKey, _ := generatePrivateKey()
		randomPrivateKeyBigInt := new(big.Int)
		randomPrivateKey.BigInt(randomPrivateKeyBigInt)

		var randomPublicKey bn254.G2Jac
		randomPublicKey.ScalarMultiplication(&g2Gen, randomPrivateKeyBigInt)
		var randomPublicKeyAffine bn254.G2Affine
		randomPublicKeyAffine.FromJacobian(&randomPublicKey)
		publicKeys = append(publicKeys, randomPublicKeyAffine)
	}
	publicKeys = append(publicKeys, rPublicKey)

	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &rPublicKey, &vPrivateKey)

	startTime := time.Now()

	for _, pk := range publicKeys {
		stealthAddress, _ := computeStealthAddress(&kPublicKey, &pk, &vPrivateKey)
		if stealthAddress.Equal(&originalStealthAddress) {
			fmt.Println("Match found!")
			/*
				missing hash
			*/
			break
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address:", duration)
}

func testSearchSpeedWithViewTag() {
	kPrivateKey, _ := generatePrivateKey()
	vPrivateKey, _ := generatePrivateKey()
	rPrivateKey, _ := generatePrivateKey()

	kPublicKey, vPublicKey, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)

	// Generate 100 random public keys in G2 and add rPublicKey as the 101st key
	_, g2Gen, _, _ := bn254.Generators()
	publicKeys := make([]bn254.G2Affine, 0, config.RunNumber+1)
	for i := 0; i < config.RunNumber; i++ {
		randomPrivateKey, _ := generatePrivateKey()
		randomPrivateKeyBigInt := new(big.Int)
		randomPrivateKey.BigInt(randomPrivateKeyBigInt)

		var randomPublicKey bn254.G2Jac
		randomPublicKey.ScalarMultiplication(&g2Gen, randomPrivateKeyBigInt)
		var randomPublicKeyAffine bn254.G2Affine
		randomPublicKeyAffine.FromJacobian(&randomPublicKey)
		publicKeys = append(publicKeys, randomPublicKeyAffine)
	}
	publicKeys = append(publicKeys, rPublicKey)

	// Compute the original stealth address
	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &rPublicKey, &vPrivateKey)

	// Calculate the view tag
	viewTag := calculateViewTag(&rPrivateKey, &vPublicKey)

	startTime := time.Now()

	// Iterate through all keys to find a match using the view tag
	for _, pk := range publicKeys {
		viewTagCalculated := calculateViewTag(&rPrivateKey, &pk)
		if viewTag == viewTagCalculated {
			temporaryStealthAddress, _ := computeStealthAddress(&kPublicKey, &pk, &vPrivateKey)
			if temporaryStealthAddress.Equal(&originalStealthAddress) {
				fmt.Println("Match found!")
				break
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address using view tag:", duration)
}

// generatePrivateKey generates a private key as a random scalar in the field.
func generatePrivateKey() (fr.Element, error) {
	var privateKey fr.Element
	_, err := privateKey.SetRandom()
	if err != nil {
		return fr.Element{}, fmt.Errorf("error generating private key: %w", err)
	}
	return privateKey, nil
}

func generatePublicKeys(kPrivateKey, vPrivateKey, rPrivateKey *fr.Element) (bn254.G1Affine, bn254.G2Affine, bn254.G2Affine) {
	// Convert private keys from fr.Element to *big.Int
	kPrivateKeyBigInt, vPrivateKeyBigInt, rPrivateKeyBigInt := new(big.Int), new(big.Int), new(big.Int)
	kPrivateKey.BigInt(kPrivateKeyBigInt)
	vPrivateKey.BigInt(vPrivateKeyBigInt)
	rPrivateKey.BigInt(rPrivateKeyBigInt)

	// Get the generators
	g1Gen, g2Gen, _, _ := bn254.Generators()

	// Scalar multiplication to generate public keys
	var kPublicKeyJac bn254.G1Jac
	var vPublicKeyJac, rPublicKeyJac bn254.G2Jac
	kPublicKeyJac.ScalarMultiplication(&g1Gen, kPrivateKeyBigInt)
	vPublicKeyJac.ScalarMultiplication(&g2Gen, vPrivateKeyBigInt)
	rPublicKeyJac.ScalarMultiplication(&g2Gen, rPrivateKeyBigInt)

	// Convert Jacobian to Affine
	var kPublicKeyAff bn254.G1Affine
	kPublicKeyAff.FromJacobian(&kPublicKeyJac)

	var vPublicKeyAff, rPublicKeyAff bn254.G2Affine
	vPublicKeyAff.FromJacobian(&vPublicKeyJac)
	rPublicKeyAff.FromJacobian(&rPublicKeyJac)

	return kPublicKeyAff, vPublicKeyAff, rPublicKeyAff
}

// computeStealthAddress computes the stealth address using pairings.
func computeStealthAddress(kPublicKey *bn254.G1Affine, rPublicKey *bn254.G2Affine, vPrivateKey *fr.Element) (bn254.GT, error) {
	// Convert vPrivateKey to big.Int for cyclotomic exponentiation
	vPrivateKeyBigInt := new(big.Int)
	vPrivateKey.BigInt(vPrivateKeyBigInt)
	// Compute pairing
	pairingResult, err := bn254.Pair([]bn254.G1Affine{*kPublicKey}, []bn254.G2Affine{*rPublicKey})
	if err != nil {
		return bn254.GT{}, fmt.Errorf("error computing pairing: %w", err)
	}

	// Compute cyclotomic exponentiation
	var stealthAddress bn254.GT
	stealthAddress.CyclotomicExp(pairingResult, vPrivateKeyBigInt)

	return stealthAddress, nil
}

func calculateViewTag(rPrivateKey *fr.Element, vPublicKey *bn254.G2Affine) uint8 {
	// Convert rPrivateKey to big.Int
	rPrivateKeyBigInt := new(big.Int)
	rPrivateKey.BigInt(rPrivateKeyBigInt)

	// Perform scalar multiplication of vPublicKey by rPrivateKey
	var vPublicKeyJac bn254.G2Jac
	var product bn254.G2Jac
	product.ScalarMultiplication(vPublicKeyJac.FromAffine(vPublicKey), rPrivateKeyBigInt)

	// Convert the product to compressed bytes
	var productAffine bn254.G2Affine
	compressedBytes := productAffine.FromJacobian(&product).Bytes()

	// Convert [64]byte array to slice
	compressedBytesSlice := compressedBytes[:]

	// Hash the compressed bytes to a field element
	domainSeparator := []byte("view_tag_domain") // Use an appropriate domain separator
	hashedFieldElements, err := fr.Hash(compressedBytesSlice, domainSeparator, 1)
	if err != nil {
		panic(fmt.Sprintf("Failed to hash to field: %v", err))
	}
	// Extract the first byte of the hashed field element as the view tag
	viewTagBytes := hashedFieldElements[0].Bytes()
	viewTag := viewTagBytes[0]

	return viewTag
}

func main() {
	// Generate private keys
	kPrivateKey, err := generatePrivateKey()
	if err != nil {
		fmt.Printf("Failed to generate kPrivateKey: %v\n", err)
		return
	}
	vPrivateKey, err := generatePrivateKey()
	if err != nil {
		fmt.Printf("Failed to generate vPrivateKey: %v\n", err)
		return
	}
	rPrivateKey, err := generatePrivateKey()
	if err != nil {
		fmt.Printf("Failed to generate rPrivateKey: %v\n", err)
		return
	}

	// Generate public keys
	kPublicKey, vPublicKey, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)
	fmt.Println("kPublicKey:", kPublicKey.Bytes())

	// Compute stealth address
	stealthAddress, err := computeStealthAddress(&kPublicKey, &rPublicKey, &vPrivateKey)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("Stealth Address:", stealthAddress)

	viewTag := calculateViewTag(&rPrivateKey, &vPublicKey)
	fmt.Println("View Tag:", viewTag)

	testSearchSpeed()
	testSearchSpeedWithViewTag()
}
