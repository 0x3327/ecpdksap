package main

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_SameP(t *testing.T) {

	// ------------------------ Generate key pairs ------------------------
	
	// ---- Recipient
	k, K, _ := utils.GenG2KeyPair()
	fmt.Println("k:", k)
	fmt.Println("K:", K)

	v, V, _ := utils.GenG1KeyPair()
	fmt.Println("v:", v)
	fmt.Println("V:", V)

	// ---- Sender
	r, R, _ := utils.GenG1KeyPair()
	fmt.Println("r:", r)
	fmt.Println("R:", R)

	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	// ---- Sender
	senderP, err := utils.SenderComputesStealthPubKey(&r, &V, &K)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("senderP: ", senderP)

	// ---- Recipient ( also Viewer )
	recipientP, err := utils.RecipientComputesStealthPubKey(&K, &R, &v)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("recipientP: ", recipientP)

	// ------------------------ -------------------------- -------------

	if senderP != recipientP {
        t.Fatalf(`ERR: sender and recipient calculated different 'P' !!!`)
    }


	str := hex.EncodeToString(k.Marshal())
	fmt.Println("k:", str)

	str = hex.EncodeToString(K.Marshal())
	fmt.Println("K:", str)

	str = hex.EncodeToString(v.Marshal())
	fmt.Println("v:", str)

	str = hex.EncodeToString(V.Marshal())
	fmt.Println("V:", str)

	str = hex.EncodeToString(r.Marshal())
	fmt.Println("r:", str)

	str = hex.EncodeToString(R.Marshal())
	fmt.Println("R:", str)
}

/*
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

*/