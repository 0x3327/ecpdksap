package sender

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func Send(jsonInputString string) (rr string, rR string, rVTag uint8, rP string) {

	// ------------------------ Unpacking json

	var senderInputData SenderInputData
	json.Unmarshal([]byte(jsonInputString), &senderInputData)

	var r fr.Element
	rBytes, _ := hex.DecodeString(senderInputData.PK_r)
	r.Unmarshal(rBytes)

	var K bn254.G2Affine
	KBytes, _ := hex.DecodeString(senderInputData.K)
	K.Unmarshal(KBytes)

	var V bn254.G1Affine
	VBytes, _ := hex.DecodeString(senderInputData.V)
	V.Unmarshal(VBytes)

	// ------------------------ Stealh Pub. Key computation
	
	R, _ := utils.CalcG1PubKey(r)
	P, _ := utils.SenderComputesStealthPubKey(&r, &V, &K);

	// ------------------------ Return val. calc.

	rr = hex.EncodeToString(r.Marshal())
	rR = hex.EncodeToString(R.Marshal())
	rVTag = utils.CalculateViewTag(&r, &V)
	rP = hex.EncodeToString(P.Marshal())

	return  rr, rR, rVTag, rP
}

type SenderInputData struct {
	PK_r string `json:"r"`
	K string `json:"K"`
	V string `json:"V"`
}

