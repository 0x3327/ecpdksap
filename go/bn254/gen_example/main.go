package gen_example

import (
	"ecpdksap-bn254/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

func GenerateExample(version string) {

	if version == "v0"{

		
	}

	if version == "v2"{

		k, K := utils.GenSECP256k1G1KeyPair()
		v, V, _ := utils.GenG1KeyPair()
		r, R, _ := utils.GenG1KeyPair()

		K_asString :=K.X.String() + "." + K.Y.String()
		V_asString :=V.X.String() + "." + V.Y.String()
		R_asString :=R.X.String() + "." + R.Y.String()

		metaInfo := MetaDbg{
			PK_k: hex.EncodeToString(k.Marshal()),
			PK_v: hex.EncodeToString(v.Marshal()),
			PK_r: hex.EncodeToString(r.Marshal()),
	
			K: K_asString,
			V: V_asString,
			R: R_asString,
	
			P_Sender: "TODO",
			VTag:     "TODO",
	
			P_Recipient: "TODO",

			Version: version,
		}

		sendParams := SendParams{
			PK_r:    metaInfo.PK_r,
			K:       metaInfo.K,
			V:       metaInfo.V,
			Version: version,
		}

		Rs, _ := utils.GenRandomRs(10)
		//add the needed `R` & its tag at the end of total array of `Rs`
		Rs = append(Rs, metaInfo.R)


		recipientParams := RecipientParams{
			PK_k:    metaInfo.PK_k,
			PK_v:    metaInfo.PK_v,
			Rs:      Rs,
			Version: version,
			WithViewTag: false, // TODO: from arg determine
		}

		pathPrefix := "./gen_example/example"

		file, _ := json.MarshalIndent(metaInfo, "", " ")
		os.WriteFile(pathPrefix + "/meta-dbg.json", file, 0644)

		file, _ = json.MarshalIndent(sendParams, "", " ")
		os.WriteFile(pathPrefix + "/inputs/send.json", file, 0644)

		file, _ = json.MarshalIndent(recipientParams, "", " ")
		os.WriteFile(pathPrefix + "/inputs/receive.json", file, 0644)

		fmt.Println("V2 executed!")
	}

	fmt.Println("Generation of example done!")
}

type MetaDbg struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	PK_r string `json:"r"`

	K string
	V string
	R string

	P_Sender string
	VTag     string

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
	VTags [] string

	Version string
	WithViewTag bool
}

