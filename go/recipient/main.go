package recipient

import (
	"ecpdksap-bn254/utils"
	ecpdksap_v0 "ecpdksap-bn254/versions/v0"
	ecpdksap_v2 "ecpdksap-bn254/versions/v2"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	BN254 "github.com/consensys/gnark-crypto/ecc/bn254"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func Scan(jsonInputString string) (rP []string, rAddr []string, privKeys []string) {

	// ------------------------ Unpacking json

	var recipientInputData RecipientInputData
	json.Unmarshal([]byte(jsonInputString), &recipientInputData)

	Rs_string := recipientInputData.Rs

	var Rs []BN254.G1Affine

	for i := 0; i < len(Rs_string); i++ {
		
		RsiX, RsiY := utils.UnpackXY(Rs_string[i])

		var Rsi BN254.G1Affine
		Rsi.X.SetString(RsiX)
		Rsi.Y.SetString(RsiY)

		Rs = append(Rs, Rsi)
	}

	if recipientInputData.Version == "v0" {

		var k, v BN254_fr.Element
		kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
		k.Unmarshal(kBytes)
		vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
		v.Unmarshal(vBytes)
	
		K, _ := utils.CalcG2PubKey(k)

		startTime := time.Now()

		if recipientInputData.WithViewTag == false { 
	
			for _, Rsi := range Rs { 
				P, _ := ecpdksap_v0.RecipientComputesStealthPubKey(&K, &Rsi, &v);
				rP = append(rP, hex.EncodeToString(P.Marshal()))
			} 
			
		} else {

			for i, Rsi := range Rs { 
				
				vR := utils.BN254_MulG1PointandElement(Rsi, v)

				if utils.BN254_G1PointToViewTag(vR, 1) != recipientInputData.ViewTags[i] { continue }

				P, _ := ecpdksap_v0.RecipientComputesStealthPubKey(&K, &Rsi, &v);
				rP = append(rP, hex.EncodeToString(P.Marshal()))
			} 
		}


		duration := time.Since(startTime)

		fmt.Println("ECPDKSAP_v2 ::: WithViewTag:", recipientInputData.WithViewTag, "; time:", duration)
	
		return 
	}

	if recipientInputData.Version == "v2" {

		var k SECP256K1_fr.Element
		kBytes, _ := hex.DecodeString(recipientInputData.PK_k)
		k.Unmarshal(kBytes)
		var k_asBigInt  big.Int
		k.BigInt(&k_asBigInt)

		var v BN254_fr.Element
		vBytes, _ := hex.DecodeString(recipientInputData.PK_v)
		v.Unmarshal(vBytes)
		var v_asBigInt  big.Int
		v.BigInt(&v_asBigInt)

		var K SECP256K1.G1Affine
		K.ScalarMultiplicationBase(&k_asBigInt)

		var V BN254.G1Affine
		V.ScalarMultiplicationBase(&v_asBigInt)

		

		_, _, _, G2_BN254 := bn254.Generators()

		var vR bn254.G1Affine
		var b SECP256K1_fr.Element
		var P SECP256K1.G1Affine
		var kb SECP256K1_fr.Element

		startTime := time.Now()

		if recipientInputData.WithViewTag == false {

			for i := 0; i < len(Rs); i++ {

				vR.ScalarMultiplication(&Rs[i], &v_asBigInt)
				S, _ := bn254.Pair([]bn254.G1Affine{vR}, []bn254.G2Affine{G2_BN254})
				b_asBigInt := ecpdksap_v2.Compute_b(&S)

				b.SetBigInt(&b_asBigInt)
				kb.Mul(&k, &b)

				P.ScalarMultiplication(&K, &b_asBigInt)

				rAddr = append(rAddr, ecpdksap_v2.ComputeEthAddress(&P))
				privKeys = append(privKeys, "0x"    + kb.Text(16))
			}
		
		} else {

			for i := 0; i < len(Rs); i++ {

				vR.ScalarMultiplication(&Rs[i], &v_asBigInt)

				if utils.BN254_G1PointToViewTag(vR, 1) != recipientInputData.ViewTags[i] {
					continue
				}
				
				S, _ := bn254.Pair([]bn254.G1Affine{vR}, []bn254.G2Affine{G2_BN254})
				b_asBigInt := ecpdksap_v2.Compute_b(&S)

				b.SetBigInt(&b_asBigInt)
				kb.Mul(&k, &b)

				P.ScalarMultiplication(&K, &b_asBigInt)

				rAddr = append(rAddr, ecpdksap_v2.ComputeEthAddress(&P))
				privKeys = append(privKeys, "0x" + kb.Text(16))
			}

		}

		duration := time.Since(startTime)

		fmt.Println("ECPDKSAP_v2 ::: WithViewTag:", recipientInputData.WithViewTag, "; time:", duration)
	}

	return
}
	
type RecipientInputData struct {
	PK_k string `json:"k"`
	PK_v string `json:"v"`
	Rs []string `json:"Rs"`
	Version string 
	ViewTags [] string
	WithViewTag bool
}