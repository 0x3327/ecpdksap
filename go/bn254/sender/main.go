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
		fmt.Printf("\nERR: 'sender.send' receives json object string (1 arg).\n\n")
		return
	}

	jsonInputString := args[2]
	var senderInputData SenderInputData
	json.Unmarshal([]byte(jsonInputString), &senderInputData)

	// ------------------------ Generate key pairs ----------------------
	
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

	fmt.Println("jsonInputString:", jsonInputString)
	fmt.Println("senderInputData:", senderInputData)
	fmt.Println("K:", K)

	P, _ := utils.SenderComputesStealthPubKey(&r, &V, &K);

	fmt.Println("-----: Sender Done!")
	fmt.Println("--- Computed P:", P)
	// ------------------------ ---------------- ------------------------

}

type SenderInputData struct {
	K string `json:"K"`
	V string `json:"V"`
}

