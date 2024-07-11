package main

import (
	"ecpdksap-bn254/utils"
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
}
