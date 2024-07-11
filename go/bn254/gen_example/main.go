package gen_example

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"os"
)


func GenerateExample () {

	// ------------------------ Generate key pairs ------------------------
	
	// ---- Recipient
	k, K, _ := utils.GenG2KeyPair()
	v, V, _ := utils.GenG1KeyPair()

	// ---- Sender
	r, R, _ := utils.GenG1KeyPair()

	// ------------------------ ---------------- ------------------------

	// ------------------------ Stealh Pub. Key computation -------------

	senderP, _ := utils.SenderComputesStealthPubKey(&r, &V, &K)

	recipientP, _ := utils.RecipientComputesStealthPubKey(&K, &R, &v)

	vTag := utils.CalculateViewTag(&r, &V)

	metaInfo := MetaDbg {
		PK_k : hex.EncodeToString(k.Marshal()),
		PK_v : hex.EncodeToString(v.Marshal()),
		PK_r : hex.EncodeToString(r.Marshal()),

		K : hex.EncodeToString(K.Marshal()),
		V : hex.EncodeToString(V.Marshal()),
		R : hex.EncodeToString(R.Marshal()),

		P_Sender: hex.EncodeToString(senderP.Marshal()),
		VTag: vTag,

		P_Recipient:  hex.EncodeToString(recipientP.Marshal()),
	}
	file, _ := json.MarshalIndent(metaInfo, "", " ")
	os.WriteFile("./cli/ex/meta-dbg.json", file, 0644)

	sendParams := SendParams {
		PK_r : metaInfo.PK_r,
		K : metaInfo.K,
		V : metaInfo.V,
	}
	file, _ = json.MarshalIndent(sendParams, "", " ")
	os.WriteFile("./cli/ex/inputs/send.json", file, 0644)

	Rs, vTags := utils.GenRandomRs(10)
	//add the needed `R` & its tag at the end of total array of `Rs`
	Rs = append(Rs, metaInfo.R)
	vTags = append(vTags, vTag)

	recipientParams := RecipientParams {
		PK_k : metaInfo.PK_k,
		PK_v : metaInfo.PK_v,
		Rs : Rs,
	}
	file, _ = json.MarshalIndent(recipientParams, "", " ")
	os.WriteFile("./cli/ex/inputs/receive.json", file, 0644)

	recipientParamsUsingVTags := RecipientParamsUsingVTags {
		PK_k : metaInfo.PK_k,
		PK_v : metaInfo.PK_v,
		Rs : Rs,
		VTags:vTags,
	}
	file, _ = json.Marshal(recipientParamsUsingVTags)
	os.WriteFile("./cli/ex/inputs/receive-using-vtag.json", file, 0644)
}

type MetaDbg struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	PK_r string `json:"r"`

	K string
	V string
	R string

	P_Sender string
	VTag uint8

	P_Recipient string
}

type SendParams struct {
	PK_r string `json:"r"`
	K string
	V string
}

type RecipientParams struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`

	Rs []string
}

type RecipientParamsUsingVTags struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`

	Rs []string
	VTags []uint8
}