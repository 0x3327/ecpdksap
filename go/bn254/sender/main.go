package sender

import (
	"ecpdksap-bn254/utils"
	ecpdksap_v2 "ecpdksap-bn254/versions/v2"
	"encoding/hex"
	"encoding/json"
	"fmt"

	BN254 "github.com/consensys/gnark-crypto/ecc/bn254"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func Send(jsonInputString string) (rr string, rR string, rVTag uint8, rP string) {

	// ------------------------ Unpacking json

	var senderInputData SenderInputData
	json.Unmarshal([]byte(jsonInputString), &senderInputData)

	fmt.Println(jsonInputString, senderInputData.Version)

	if senderInputData.Version == "v0" {

		var r BN254_fr.Element
		rBytes, _ := hex.DecodeString(senderInputData.PK_r)
		r.Unmarshal(rBytes)

		var K BN254.G2Affine
		KBytes, _ := hex.DecodeString(senderInputData.K)
		K.Unmarshal(KBytes)

		var V BN254.G1Affine
		VBytes, _ := hex.DecodeString(senderInputData.V)
		V.Unmarshal(VBytes)

		// ------------------------ Stealh Pub. Key computation

		R, _ := utils.CalcG1PubKey(r)
		P, _ := utils.SenderComputesStealthPubKey(&r, &V, &K)

		// ------------------------ Return val. calc.

		rr = hex.EncodeToString(r.Marshal())
		rR = hex.EncodeToString(R.Marshal())
		rVTag = utils.CalculateViewTag(&r, &V)
		rP = hex.EncodeToString(P.Marshal())

	}

	if senderInputData.Version == "v2" {

		var K SECP256K1.G1Affine
		Kx, Ky := utils.UnpackXY(senderInputData.K)
		K.X.SetString(Kx)
		K.Y.SetString(Ky)

		var V BN254.G1Affine
		Vx, Vy := utils.UnpackXY(senderInputData.V)
		V.X.SetString(Vx)
		V.Y.SetString(Vy)

		var r BN254_fr.Element
		rBytes, _ := hex.DecodeString(senderInputData.PK_r)
		r.Unmarshal(rBytes)

		GT, _ := ecpdksap_v2.SenderComputesSharedSecret(&r, &V, &K)

		b := ecpdksap_v2.Compute_b(&GT)
		var b_asElement SECP256K1_fr.Element
		b_asElement.SetBigInt(&b)

		ethAddr := ecpdksap_v2.SenderComputesEthAddress(&b_asElement, &K)

		fmt.Println("V2 executed :: ethAddr:", ethAddr)
	}

	return rr, rR, rVTag, rP
}

type SenderInputData struct {
	PK_r    string `json:"r"`
	K       string `json:"K"`
	V       string `json:"V"`
	Version string
}
