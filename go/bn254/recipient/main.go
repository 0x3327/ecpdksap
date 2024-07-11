package recipient

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func Scan(jsonInputString string) (rP []string) {

	// ------------------------ Unpacking json

	var recipientInputData RecipientInputData
	json.Unmarshal([]byte(jsonInputString), &recipientInputData)

	var k, v fr.Element
	kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
	k.Unmarshal(kBytes)
	vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
	v.Unmarshal(vBytes)

	Rs := recipientInputData.Rs

	// ------------------------ Stealh Pub. Key computation

	K, _ := utils.CalcG2PubKey(k)
	// V, _ := utils.CalcG1PubKey(v) - not needed

	var R bn254.G1Affine
	for _, R_string := range Rs { 
		R_bytes, _ := hex.DecodeString(R_string)
		R.Unmarshal(R_bytes)
		P, _ := utils.RecipientComputesStealthPubKey(&K, &R, &v);
		rP = append(rP, hex.EncodeToString(P.Marshal()))
    } 

	return rP
}


func ScanUsingViewTag(jsonInputString string) (rP []string){

	// ------------------------ Unpacking json

	var recipientInputData RecipientInputDataWithViewTag
	json.Unmarshal([]byte(jsonInputString), &recipientInputData)

	var k, v fr.Element
	kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
	k.Unmarshal(kBytes)
	vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
	v.Unmarshal(vBytes)

	var vTags = recipientInputData.ViewTags
	Rs := recipientInputData.Rs

	// ------------------------ Stealh Pub. Key computation

	K, _ := utils.CalcG2PubKey(k)

	var R bn254.G1Affine
	for idx, R_string := range Rs { 
		R_bytes, _ := hex.DecodeString(R_string)
		R.Unmarshal(R_bytes)

		currViewTag := utils.CalculateViewTag(&v, &R);
		fmt.Println(currViewTag, vTags[idx])

		if currViewTag == vTags[idx] {
			P, _ := utils.RecipientComputesStealthPubKey(&K, &R, &v);
			rP = append(rP, hex.EncodeToString(P.Marshal()))
		}
    } 

	return rP
}

type RecipientInputData struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	Rs []string `json:"Rs"`
}

type RecipientInputDataWithViewTag struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	ViewTags []uint8 `json:"VTags"`
	Rs []string `json:"Rs"`
}