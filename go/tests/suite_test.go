package main

import (
	"math/big"
	"testing"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"

	ecpdksap_v0 "ecpdksap-go/versions/v0"
	ecpdksap_v1 "ecpdksap-go/versions/v1"
	ecpdksap_v2 "ecpdksap-go/versions/v2"

	"ecpdksap-go/utils"
)

func Test_V0(t *testing.T) {

	_, K, _ := utils.BN254_GenG2KeyPair()
	v, V, _ := utils.BN254_GenG1KeyPair()

	r, R, _ := utils.BN254_GenG1KeyPair()

	P_Sender, _ := ecpdksap_v0.SenderComputesStealthPubKey(&r, &V, &K)

	P_Recipient, _ := ecpdksap_v0.RecipientComputesStealthPubKey(&K, &R, &v)

	if P_Sender != P_Recipient {
		t.Fatalf(`ERR: sender and recipient calculated different public key !!!`)
	}
}

func Test_V1(t *testing.T) {

	k, K, _ := utils.BN254_GenG2KeyPair()
	v, V, _ := utils.BN254_GenG1KeyPair()

	r, R, _ := utils.BN254_GenG1KeyPair()

	P_Sender, _ := ecpdksap_v1.SenderComputesStealthPubKey(&r, &V, &K)

	P_Recipient := ecpdksap_v1.RecipientComputesStealthPubKey(&k, &v, &R)

	P_Viewer := ecpdksap_v1.ViewerComputesStealthPubKey(&K, &R, &v)

	if P_Sender != P_Recipient {
		t.Fatalf(`ERR: sender and recipient calculated different public key !!!`)
	}

	if P_Viewer != P_Recipient {
		t.Fatalf(`ERR: viewer calculated different public key than sender and recipient !!!`)
	}
}

func Test_V2(t *testing.T) {

	k, K := utils.SECP256k_Gen1G1KeyPair()
	v, V, _ := utils.BN254_GenG1KeyPair()

	r, R, _ := utils.BN254_GenG1KeyPair()

	S_Sender := ecpdksap_v2.SenderComputesSharedSecret(&r, &V, &K)

	S_Recipient := ecpdksap_v2.RecipientComputesSharedSecret(&v, &R, &K)

	if S_Sender != S_Recipient {
		t.Fatalf(`ERR: sender and recipient calculated different secret !!!`)
	}

	b_Sender := ecpdksap_v2.Compute_b_asElement(&S_Sender)
	b_Recipient := ecpdksap_v2.Compute_b_asElement(&S_Recipient)

	if b_Sender.String() != b_Recipient.String() {
		t.Fatalf(`ERR: sender and recipient calculated different 'b' !!!`)
	}

	ethAddr_Sender := ecpdksap_v2.SenderComputesEthAddress(&b_Sender, &K)

	var P SECP256K1.G1Affine

	_, G1 := SECP256K1.Generators()
	var kb SECP256K1_fr.Element
	kb.Mul(&k, &b_Recipient)
	P.ScalarMultiplication(&G1, kb.BigInt(new (big.Int)))
	ethAddr_Recipient := ecpdksap_v2.ComputeEthAddress(&P)

	if ethAddr_Sender!= ethAddr_Recipient {
		t.Fatalf(`ERR: sender and recipient calculated different ETH addr !!!`)
	}

}
