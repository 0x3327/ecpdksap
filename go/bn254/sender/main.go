package sender

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
)

func Send() {

	args := os.Args

	if len(args) != 3 {
		fmt.Printf("\nERR: 'sender.send' receives only 1 arg.\n\n")
		return
	}

	jsonInputString := args[2]

	byteValue, _ := hex.DecodeString(jsonInputString)

	var senderInputData SenderInputData
	json.Unmarshal(byteValue, &senderInputData)

	// ------------------------ Generate key pairs ------------------------
	
	// ---- Sender
	r, R, err := utils.GenG1KeyPair()
	if err != nil {
		fmt.Printf("Failed to generate rPrivateKey: %v\n", err)
		return
	}
	fmt.Println("r:", r)
	fmt.Println("R:", R)

	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	var K bn254.G2Affine
	var V bn254.G1Affine
	KBytes, _ := hex.DecodeString(senderInputData.K)
	K.Unmarshal(KBytes)
	VBytes, _ := hex.DecodeString(senderInputData.V)
	V.Unmarshal(VBytes)

	utils.SenderComputesStealthPubKey(&r, &V, &K);


	fmt.Println("-----> Sender Done!")

	// ------------------------ ---------------- ------------------------

}



type SenderInputData struct {
	K string `json:"K"`
	V string `json:"V"`
}

