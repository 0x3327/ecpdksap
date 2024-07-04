package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"sap-go/config"
	"time"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func testSearchSpeed() {
	kPrivateKey, _ := generatePrivateKey()
	vPrivateKey, _ := generatePrivateKey()
	rPrivateKey, _ := generatePrivateKey()

	kPublicKey, vPublicKey, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)

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

	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &vPublicKey, &rPublicKey)
	formattedOriginalStealthAddress := formatStealthAddress(&originalStealthAddress)

	startTime := time.Now()

	for _, pk := range publicKeys {
		stealthAddress, _ := computeStealthAddress(&kPublicKey, &vPublicKey, &pk)
		formattedStealthAddress := formatStealthAddress(&stealthAddress)
		if formattedStealthAddress == formattedOriginalStealthAddress {
			fmt.Println("Match found!")
			break
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address:", duration)
}

func testSearchSpeedWithViewTag() time.Duration {
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
	originalStealthAddress, _ := computeStealthAddress(&kPublicKey, &vPublicKey, &rPublicKey)
	formattedOriginalStealthAddress := formatStealthAddress(&originalStealthAddress)

	// Calculate the view tag
	viewTag, err := calculateViewTag(&rPublicKey, &vPublicKey)
	if err != nil {
		fmt.Println("Error computing view tag:", err)
		return time.Duration(0)
	}

	startTime := time.Now()

	// Iterate through all keys to find a match using the view tag
	for _, pk := range publicKeys {
		viewTagCalculated, err := calculateViewTag(&pk, &vPublicKey) // reciever uses his private viewing key and sender public key to calculate the view tag
		if err != nil {
			fmt.Println("Error computing view tag:", err)
			return time.Duration(0)
		}
		if viewTag == viewTagCalculated {
			temporaryStealthAddress, _ := computeStealthAddress(&kPublicKey, &vPublicKey, &pk)
			formattedTemporaryStealthAddress := formatStealthAddress(&temporaryStealthAddress)
			if formattedTemporaryStealthAddress == formattedOriginalStealthAddress {
				fmt.Println("Match found!")
				break
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Println("Time taken to find the address using view tag:", duration)

	return duration
}

func runExperiment() {
	results := make([][]string, 0, 12)
	results = append(results, []string{"Run", "Duration (ms)", "Public Keys"})

	var totalDuration time.Duration
	for i := 0; i < 10; i++ {
		duration := testSearchSpeedWithViewTag()
		totalDuration += duration

		results = append(results, []string{fmt.Sprintf("%d", i+1), fmt.Sprintf("%.2f", float64(duration.Milliseconds())), fmt.Sprintf("%d", config.RunNumber)})
	}

	avgDurationMs := float64(totalDuration.Milliseconds()) / 10.0

	results = append(results, []string{"Average", fmt.Sprintf("%.2f", avgDurationMs), fmt.Sprintf("%d", config.RunNumber)})

	// Save results to CSV file
	fileName := fmt.Sprintf("experiment_results_%d_public_keys.csv", config.RunNumber)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(results)
	if err != nil {
		fmt.Println("Error writing to CSV file:", err)
		return
	}

	fmt.Println("Experiment results saved to", fileName)
}

func hash(input []byte) []byte {
	hasher := sha256.New()
	hasher.Write(input)     // Hash the input
	hash := hasher.Sum(nil) // Finalize the hash and return the result
	return hash
}

func hashToField(input []byte) [32]byte {
	domainSeparator := []byte("view_tag_domain") // Use an appropriate domain separator
	hashedFieldElements, err := fr.Hash(input, domainSeparator, 1)
	if err != nil {
		panic(fmt.Sprintf("Failed to hash to field: %v", err))
	}
	// Extract the first byte of the hashed field element as the view tag
	return hashedFieldElements[0].Bytes()
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

func generatePublicKeys(kPrivateKey, vPrivateKey, rPrivateKey *fr.Element) (bn254.G1Affine, bn254.G2Affine, bn254.G1Affine) {
	// Convert private keys from fr.Element to *big.Int
	kPrivateKeyBigInt, vPrivateKeyBigInt, rPrivateKeyBigInt := new(big.Int), new(big.Int), new(big.Int)
	kPrivateKey.BigInt(kPrivateKeyBigInt)
	vPrivateKey.BigInt(vPrivateKeyBigInt)
	rPrivateKey.BigInt(rPrivateKeyBigInt)

	// Get the generators
	g1Gen, g2Gen, _, _ := bn254.Generators()

	// Scalar multiplication to generate public keys
	var kPublicKeyJac, rPublicKeyJac bn254.G1Jac
	var vPublicKeyJac bn254.G2Jac
	kPublicKeyJac.ScalarMultiplication(&g1Gen, kPrivateKeyBigInt)
	vPublicKeyJac.ScalarMultiplication(&g2Gen, vPrivateKeyBigInt)
	rPublicKeyJac.ScalarMultiplication(&g1Gen, rPrivateKeyBigInt)

	// Convert Jacobian to Affine
	var kPublicKeyAff, rPublicKeyAff bn254.G1Affine
	kPublicKeyAff.FromJacobian(&kPublicKeyJac)
	rPublicKeyAff.FromJacobian(&rPublicKeyJac)

	var vPublicKeyAff bn254.G2Affine
	vPublicKeyAff.FromJacobian(&vPublicKeyJac)

	return kPublicKeyAff, vPublicKeyAff, rPublicKeyAff
}

// computeStealthAddress computes the stealth address using pairings.
func computeStealthAddress(kPublicKey *bn254.G1Affine, vPublicKey *bn254.G2Affine, rPublicKey *bn254.G1Affine) (bn254.G1Affine, error) {
	g1Gen, _, _, _ := bn254.Generators()
	// Compute pairing
	pairingResult, err := bn254.Pair([]bn254.G1Affine{*rPublicKey}, []bn254.G2Affine{*vPublicKey})
	// fmt.Println("pairingResult in bytes:", pairingResult.Bytes())
	if err != nil {
		return bn254.G1Affine{}, fmt.Errorf("error computing pairing: %w", err)
	}
	// Compute shared secret
	pairingResultBytes := pairingResult.Bytes()
	sharedSecretHashed := hashToField(pairingResultBytes[:])
	sharedSecretBigInt := new(big.Int).SetBytes(sharedSecretHashed[:])

	// Convert shared secret to G1 point
	var sharedSecretJac bn254.G1Jac
	sharedSecretJac.ScalarMultiplication(&g1Gen, sharedSecretBigInt)

	// Compute stealth address
	var stealthAddressJac bn254.G1Jac
	stealthAddressJac.FromAffine(kPublicKey) // for efficiency, we use the Jacobian representation
	stealthAddressJac.AddAssign(&sharedSecretJac)

	var stealthAddressAff bn254.G1Affine
	stealthAddressAff.FromJacobian(&stealthAddressJac)

	return stealthAddressAff, nil
}

func formatStealthAddress(stealthAddress *bn254.G1Affine) string {
	stealthAddressBytes := stealthAddress.Bytes()
	return "0x" + fmt.Sprintf("%x", hash(stealthAddressBytes[:])[:20])
}

func calculateViewTag(rPublicKey *bn254.G1Affine, vPublicKey *bn254.G2Affine) (uint8, error) {
	// Compute pairing
	pairingResult, err := bn254.Pair([]bn254.G1Affine{*rPublicKey}, []bn254.G2Affine{*vPublicKey})
	if err != nil {
		return 0, fmt.Errorf("error computing pairing: %w", err)
	}
	// Compute view tag
	pairingResultBytes := pairingResult.Bytes()
	sharedSecretHashed := hashToField(pairingResultBytes[:])
	// Extract the first byte of the hashed field element as the view tag
	viewTag := sharedSecretHashed[0]

	return viewTag, nil
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
	vPrivateKeyBigInt := new(big.Int)
	vPrivateKey.BigInt(vPrivateKeyBigInt)
	// Generate public keys
	kPublicKey, vPublicKey, rPublicKey := generatePublicKeys(&kPrivateKey, &vPrivateKey, &rPrivateKey)
	// fmt.Println("kPublicKey:", kPublicKey)
	// fmt.Println("vPublicKey:", vPublicKey)
	// fmt.Println("rPublicKey:", rPublicKey)

	// Compute stealth address
	stealthAddress, err := computeStealthAddress(&kPublicKey, &vPublicKey, &rPublicKey)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("Field Stealth Address Representation:", stealthAddress)
	fmt.Println("Formatted Stealth Address:", formatStealthAddress(&stealthAddress))

	viewTag, err := calculateViewTag(&rPublicKey, &vPublicKey)
	if err != nil {
		fmt.Println("Error computing view tag:", err)
		return
	}
	fmt.Println("View Tag:", viewTag)

	startTime := time.Now()
	pairingResult, err := bn254.Pair([]bn254.G1Affine{rPublicKey}, []bn254.G2Affine{vPublicKey})
	// fmt.Println("pairingResult in bytes:", pairingResult.Bytes())
	if err != nil {
		fmt.Errorf("error computing pairing: %w", err)
	}
	duration := time.Since(startTime)
	fmt.Println("Time taken to compute pairing:", duration)
	startTime = time.Now()
	pairingResult, err = bn254.Pair([]bn254.G1Affine{rPublicKey}, []bn254.G2Affine{vPublicKey})
	// fmt.Println("pairingResult in bytes:", pairingResult.Bytes())
	if err != nil {
		fmt.Errorf("error computing pairing: %w", err)
	}
	pairingResult.CyclotomicExp(pairingResult, vPrivateKeyBigInt)
	duration = time.Since(startTime)
	fmt.Println("Time taken to compute pairing and cyclotomic exp:", duration)

	// testSearchSpeed()
	// testSearchSpeedWithViewTag()
	// runExperiment()
}
