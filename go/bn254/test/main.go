package main

import (
	"ecpdksap-bn254/config"
	"fmt"
	"math/big"
	"testing"
	"time"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
)


func TestSearchSpeed(t *testing.T) {
	kPrivateKey, _ := genG1KeyPair()
	vPrivateKey, _ := generatePrivateKey()
	rPrivateKey, _ := generatePrivateKey()

	kPublicKey, _, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)

	g1Gen, _, _, _ := bn254.Generators()
	publicKeys := make([]bn254.G1Affine, 0, config.RunNumber+1)
	for i := 0; i < config.RunNumber; i++ {
		randomPrivateKey, _ := generatePrivateKey()
		randomPrivateKeyBigInt := new(big.Int)
		randomPrivateKey.BigInt(randomPrivateKeyBigInt)

		var randomPublicKey bn254.G1Jac
		randomPublicKey.ScalarMultiplication(&g1Gen, randomPrivateKeyBigInt)
		var randomPublicKeyAffine bn254.G1Affine
		randomPublicKeyAffine.FromJacobian(&randomPublicKey)
		publicKeys = append(publicKeys, randomPublicKeyAffine)
	}
	publicKeys = append(publicKeys, rPublicKey)

	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &rPublicKey, &vPrivateKey)
	formattedOriginalStealthAddress := formatStealthAddress(&originalStealthAddress)

	startTime := time.Now()

	for _, pk := range publicKeys {
		stealthAddress, _ := computeStealthAddress(&kPublicKey, &pk, &vPrivateKey)
		formattedStealthAddress := formatStealthAddress(&stealthAddress)
		if formattedStealthAddress == formattedOriginalStealthAddress {
			fmt.Println("Match found!")
			break
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address:", duration)
}



func TestSearchSpeedWithViewTag(t *testing.T)  {
	kPrivateKey, _ := generatePrivateKey()
	vPrivateKey, _ := generatePrivateKey()
	rPrivateKey, _ := generatePrivateKey()

	kPublicKey, vPublicKey, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)

	// Generate 100 random public keys in G2 and add rPublicKey as the 101st key
	g1Gen, _, _, _ := bn254.Generators()
	publicKeys := make([]bn254.G1Affine, 0, config.RunNumber+1)
	for i := 0; i < config.RunNumber; i++ {
		randomPrivateKey, _ := generatePrivateKey()
		randomPrivateKeyBigInt := new(big.Int)
		randomPrivateKey.BigInt(randomPrivateKeyBigInt)

		var randomPublicKey bn254.G1Jac
		randomPublicKey.ScalarMultiplication(&g1Gen, randomPrivateKeyBigInt)
		var randomPublicKeyAffine bn254.G1Affine
		randomPublicKeyAffine.FromJacobian(&randomPublicKey)
		publicKeys = append(publicKeys, randomPublicKeyAffine)
	}
	publicKeys = append(publicKeys, rPublicKey)

	// Compute the original stealth address
	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &rPublicKey, &vPrivateKey)
	formattedOriginalStealthAddress := formatStealthAddress(&originalStealthAddress)

	// Calculate the view tag
	viewTag := calculateViewTag(&rPrivateKey, &vPublicKey)

	startTime := time.Now()

	// Iterate through all keys to find a match using the view tag
	for _, pk := range publicKeys {
		viewTagCalculated := calculateViewTag(&vPrivateKey, &pk) // reciever uses his private viewing key and sender public key to calculate the view tag
		if viewTag == viewTagCalculated {
			temporaryStealthAddress, _ := computeStealthAddress(&kPublicKey, &pk, &vPrivateKey)
			formattedTemporaryStealthAddress := formatStealthAddress(&temporaryStealthAddress)
			if formattedTemporaryStealthAddress == formattedOriginalStealthAddress {
				fmt.Println("Match found!")
				break
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address using view tag:", duration)

}