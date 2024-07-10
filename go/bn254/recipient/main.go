package recipient

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func Scan() {

	args := os.Args

	if len(args) != 3 {
		fmt.Printf("\nERR: 'recipient.scan' receives json object string (1 arg).\n\n")
		return
	}

	jsonInputString := args[2]
	var recipientInputData RecipientInputData
	json.Unmarshal([]byte(jsonInputString), &recipientInputData)

	// ------------------------ Stealh Pub. Key computation -------------

	var k, v fr.Element
	kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
	k.Unmarshal(kBytes)
	vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
	v.Unmarshal(vBytes)

	fmt.Println(k)

	Rs := recipientInputData.Rs

	fmt.Println(recipientInputData.Rs)

	K, _ := utils.CalcG2PubKey(k)
	// V, _ := utils.CalcG1PubKey(v)

	var R bn254.G1Affine
	for _, R_string := range Rs { 
		R_bytes, _ := hex.DecodeString(R_string)
		R.Unmarshal(R_bytes)
		P, _ := utils.RecipientComputesStealthPubKey(&K, &R, &v);

		fmt.Println("computed P: ", P)
    } 

	fmt.Println("-----: Recipient Done!")
	// ------------------------ ---------------- ------------------------
}


func ScanUsingViewTag() {

	args := os.Args

	if len(args) != 3 {
		fmt.Printf("\nERR: 'recipient.scan-with-vtag' receives json object string (1 arg).\n\n")
		return
	}

	jsonInputString := args[2]
	var recipientInputData RecipientInputDataWithViewTag
	json.Unmarshal([]byte(jsonInputString), &recipientInputData)

	// ------------------------ Stealh Pub. Key computation -------------

	var k, v fr.Element
	kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
	k.Unmarshal(kBytes)
	vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
	v.Unmarshal(vBytes)

	var vTags = recipientInputData.ViewTags
	Rs := recipientInputData.Rs

	fmt.Println(vTags, Rs)

	K, _ := utils.CalcG2PubKey(k)

	var R bn254.G1Affine
	for idx, R_string := range Rs { 
		R_bytes, _ := hex.DecodeString(R_string)
		R.Unmarshal(R_bytes)

		currViewTag := utils.CalculateViewTag(&v, &R);
		if currViewTag == vTags[idx] {
			P, _ := utils.RecipientComputesStealthPubKey(&K, &R, &v);
			fmt.Println("computed P: ", P)
		}
    } 

	fmt.Println("-----: Recipient Done!")
	// ------------------------ ---------------- ------------------------
}

type RecipientInputData struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	Rs []string `json:"Rs"`
}

type RecipientInputDataWithViewTag struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	ViewTags []uint8 `json:"viewTags"`
	Rs []string `json:"Rs"`
}