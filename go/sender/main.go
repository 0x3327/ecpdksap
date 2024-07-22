package sender

import (
	"encoding/hex"
	"encoding/json"

	BN254 "github.com/consensys/gnark-crypto/ecc/bn254"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"

	ecpdksap_v0 "ecpdksap-go/versions/v0"
	ecpdksap_v1 "ecpdksap-go/versions/v1"
	ecpdksap_v2 "ecpdksap-go/versions/v2"

	"ecpdksap-go/utils"
)

func Send(jsonInputString string) (rr string, rR string, rVTag uint8, rP string) {

	// ------------------------ Unpacking json

	var senderInputData SenderInputData
	json.Unmarshal([]byte(jsonInputString), &senderInputData)

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

		R, _ := utils.BN254_CalcG1PubKey(r)
		P, _ := ecpdksap_v0.SenderComputesStealthPubKey(&r, &V, &K)

		// ------------------------ Return val. calc.

		rr = hex.EncodeToString(r.Marshal())
		rR = hex.EncodeToString(R.Marshal())
		rVTag = ecpdksap_v0.CalculateViewTag(&r, &V)
		rP = hex.EncodeToString(P.Marshal())

	} else if senderInputData.Version == "v1" {

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

		R, _ := utils.BN254_CalcG1PubKey(r)
		P, _ := ecpdksap_v1.SenderComputesStealthPubKey(&r, &V, &K)

		// ------------------------ Return val. calc.

		rr = hex.EncodeToString(r.Marshal())
		rR = hex.EncodeToString(R.Marshal())
		rVTag = ecpdksap_v1.CalculateViewTag(&r, &V)
		rP = hex.EncodeToString(P.Marshal())

	} else if senderInputData.Version == "v2" {

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

		GT := ecpdksap_v2.SenderComputesSharedSecret(&r, &V, &K)

		b := ecpdksap_v2.Compute_b(&GT)
		var b_asElement SECP256K1_fr.Element
		b_asElement.SetBigInt(&b)

		ecpdksap_v2.SenderComputesEthAddress(&b_asElement, &K)
	}

	return rr, rR, rVTag, rP
}

type SenderInputData struct {
	PK_r    string `json:"r"`
	K       string `json:"K"`
	V       string `json:"V"`
	Version string
}
