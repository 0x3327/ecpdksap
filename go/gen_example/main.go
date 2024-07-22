package gen_example

import (
	"ecpdksap-go/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func GenerateExample(version string, sampleSizeStr string) (sendParams SendParams, recipientParams RecipientParams) {

	v, V, _ := utils.BN254_GenG1KeyPair()
	r, R, _ := utils.BN254_GenG1KeyPair()

	var K_asString string
	var kBytes []byte

	if version == "v0" {
		k, K, _ := utils.BN254_GenG2KeyPair()
		K_asString = K.X.String() + "." + K.Y.String()
		kBytes = k.Marshal()
	}

	if version == "v2" {
		k, K := utils.SECP256k_Gen1G1KeyPair()
		K_asString = K.X.String() + "." + K.Y.String()
		kBytes = k.Marshal()
	}

	V_asString := V.X.String() + "." + V.Y.String()
	R_asString := R.X.String() + "." + R.Y.String()

	tmp := utils.BN254_MulG1PointandElement(V, r)
	viewTag := utils.BN254_G1PointToViewTag(&tmp, 1)

	metaInfo := MetaDbg{
		PK_k: hex.EncodeToString(kBytes),
		PK_v: hex.EncodeToString(v.Marshal()),
		PK_r: hex.EncodeToString(r.Marshal()),

		K: K_asString,
		V: V_asString,
		R: R_asString,

		ViewTag:  viewTag,
		Version: version,

		P_Sender: "TODO",
		P_Recipient: "TODO",
	}

	sendParams = SendParams{
		PK_r:    metaInfo.PK_r,
		K:       metaInfo.K,
		V:       metaInfo.V,
		Version: version,
	}

	sampleSize, _ := strconv.Atoi(sampleSizeStr)
	Rs, viewTags := utils.GenRandomRsAndViewTags(sampleSize)
	Rs = append(Rs, metaInfo.R)
	viewTags = append(viewTags, metaInfo.ViewTag)


	recipientParams = RecipientParams{
		PK_k:    metaInfo.PK_k,
		PK_v:    metaInfo.PK_v,
		Rs:      Rs,
		Version: version,
		WithViewTag: false, // TODO: from arg determine
		ViewTags: viewTags,
	}

	pathPrefix := "./gen_example/example"

	file, _ := json.MarshalIndent(metaInfo, "", " ")
	os.WriteFile(pathPrefix + "/meta-dbg.json", file, 0644)

	file, _ = json.MarshalIndent(sendParams, "", " ")
	os.WriteFile(pathPrefix + "/inputs/send.json", file, 0644)

	file, _ = json.MarshalIndent(recipientParams, "", " ")
	os.WriteFile(pathPrefix + "/inputs/receive.json", file, 0644)

	fmt.Println("Example generation for", version, "done!")

	return
}

type MetaDbg struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	PK_r string `json:"r"`

	K string
	V string
	R string

	P_Sender string
	ViewTag     string

	P_Recipient string

	Version string
}

type SendParams struct {
	PK_r string `json:"r"`
	K    string
	V    string

	Version string
}

type RecipientParams struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`

	Rs []string
	ViewTags [] string

	Version string
	WithViewTag bool
}

