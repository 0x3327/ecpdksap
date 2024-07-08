package main

import (
	"ecpdksap-bn254/sender"
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"fmt"
	"os"
)

func MainSend (jsonInputString string) {

	fmt.Println(jsonInputString)
}


func main () {

	args := os.Args

	if len(args) == 1 {
		fmt.Printf("\nERR: no args passed!\n\n")
		return
	}

	subcmd := args[1];

	if subcmd != "send" && subcmd != "receive-scan" {
		fmt.Printf("\nERR: only 'send' and 'receive-scan' subcommands allowed.\n\n")
		return
	}

	if subcmd == "send" {
		sender.Send()
	}

	k, K, _ := utils.GenG2KeyPair()
	v, V, _ := utils.GenG1KeyPair()
	fmt.Println("k:", k)
	fmt.Println("K:", K.Marshal())
	fmt.Println("v:", v)
	fmt.Println("V:", V.Marshal())


	str := hex.EncodeToString(K.Marshal())
	fmt.Println("K:", str)
	str = hex.EncodeToString(V.Marshal())
	fmt.Println("V:", str)
	// b, _ := hex.DecodeString(str)
	// fmt.Println(b)
}



func main2() {
	// ------------------------ Generate key pairs ------------------------
	
	// ---- Recipient
	k, K, err := utils.GenG2KeyPair()
	if err != nil {
		fmt.Printf("Failed to generate kPrivateKey: %v\n", err)
		return
	}
	fmt.Println("k:", k)
	fmt.Println("K:", K)
	v, V, err := utils.GenG1KeyPair()
	if err != nil {
		fmt.Printf("Failed to generate vPrivateKey: %v\n", err)
		return
	}
	fmt.Println("v:", v)
	fmt.Println("V:", V)
	// ---- Sender
	r, R, err := utils.GenG1KeyPair()
	if err != nil {
		fmt.Printf("Failed to generate rPrivateKey: %v\n", err)
		return
	}
	fmt.Println("r:", r)

	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	// ---- Sender


	senderP, err := utils.SenderComputesStealthPubKey(&r, &V, &K)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("senderP: ", senderP)

	// ---- Recipient
	recipientP, err := utils.RecipientComputesStealthPubKey(&K, &R, &v)
	if err != nil {
		fmt.Println("Error computing stealth address:", err)
		return
	}
	fmt.Println("recipientP: ", recipientP)

	// recipient2P, err := computeStealthAddress(&K, &R, &v)
	// fmt.Println("recipient2P: ", recipient2P)

	viewTag := utils.CalculateViewTag(&r, &V)
	fmt.Println("View Tag:", viewTag)


	// ------------------------ -------------------------- -------------

}

