package main

import (
	"ecpdksap-bn254/utils"
	"fmt"
	"testing"
)

func Test_V0_SameP(t *testing.T) {

	// ------------------------ Generate key pairs ------------------------

	// ---- Recipient
	_, K, _ := utils.GenG2KeyPair()
	v, V, _ := utils.GenG1KeyPair()

	// ---- Sender
	r, R, _ := utils.GenG1KeyPair()

	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	// ---- Sender
	senderP, err := utils.SenderComputesStealthPubKey(&r, &V, &K)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}

	// ---- Recipient ( also Viewer )
	recipientP, err := utils.RecipientComputesStealthPubKey(&K, &R, &v)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}

	// ------------------------ -------------------------- -------------

	if senderP != recipientP {
		t.Fatalf(`ERR: sender and recipient calculated different 'P' !!!`)
	}
}

// func Test_V1_SameP(t *testing.T) {

// 	// ------------------------ Generate key pairs ------------------------

// 	// ---- Recipient
// 	_, K, _ := utils.GenG2KeyPair()
// 	v, V, _ := utils.GenG1KeyPair()

// 	// ---- Sender
// 	r, R, _ := utils.GenG1KeyPair()

// 	// ------------------------ ---------------- ------------------------
// 	// sender =======> recipient
// 	// e(R * hash(r*V), K) = e(R * hash(R*v), K)

// 	// utils.CalcG1PubKey()

// 	// g1Gen, _, _, _ := bn254.Generators()

// 	rBigInt := new(big.Int)
// 	r.BigInt(rBigInt)
// 	var RV_product bn254.G1Jac
// 	var VJac bn254.G1Jac
// 	VJac.FromAffine(&V)
// 	RV_product.ScalarMultiplication(&VJac, rBigInt)

// 	var RV_productAff bn254.G1Affine
// 	RV_productAff.FromJacobian(&RV_product)

// 	RV_productTemp := RV_productAff.Bytes()

// 	var h fr.Element
// 	h.SetBytes(utils.Hash(RV_productTemp[:]))
// 	hBigInt := new(big.Int)
// 	h.BigInt(hBigInt)

// 	var new_r fr.Element
// 	new_r.SetBytes(utils.Hash(KBytes[:]))

// 	// ------------------------ Stealh Pub. Key computation -------------

// 	// ---- Sender
// 	senderP, err := utils.SenderComputesStealthPubKey(&r, &V, &K)
// 	if err != nil {
// 		fmt.Println("Error computing stealth address:", err)
// 		return
// 	}

// 	// ---- Recipient ( also Viewer )
// 	recipientP, err := utils.RecipientComputesStealthPubKey(&K, &R, &v)
// 	if err != nil {
// 		fmt.Println("Error computing stealth address:", err)
// 		return
// 	}

// 	// ------------------------ -------------------------- -------------

// 	if senderP != recipientP {
//         t.Fatalf(`ERR: sender and recipient calculated different 'P' !!!`)
//     }
// }
