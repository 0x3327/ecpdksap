package sender

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
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
	


	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	var r fr.Element
	rBytes, _ := hex.DecodeString(senderInputData.PK_r)
	r.Unmarshal(rBytes)
	R, _ := utils.CalcG1PubKey(r)

	var K bn254.G2Affine
	var V bn254.G1Affine

	KBytes, _ := hex.DecodeString(senderInputData.K)
	K.Unmarshal(KBytes)
	VBytes, _ := hex.DecodeString(senderInputData.V)
	V.Unmarshal(VBytes)

	P, _ := utils.SenderComputesStealthPubKey(&r, &V, &K);

	fmt.Println("-----: Sender Done! :-----")

	fmt.Println("\n- DBG -")
	fmt.Println("P:", hex.EncodeToString(P.Marshal()))
	fmt.Println("r:", r)

	fmt.Println("\n- Public info -\n")
	fmt.Println("R:", hex.EncodeToString(R.Marshal()))
	fmt.Println()
	vTag := utils.CalculateViewTag(&r, &V)
	fmt.Println("vTag:", vTag)
	fmt.Println()

	// ------------------------ ---------------- ------------------------

}

type SenderInputData struct {
	PK_r string `json:"r"`
	K string `json:"K"`
	V string `json:"V"`
}

