package main

import (
	"ecpdksap-go/utils"
	ecpdksap_v0 "ecpdksap-go/versions/v0"
	ecpdksap_v2 "ecpdksap-go/versions/v2"
	"testing"
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

func Test_V2(t *testing.T) {

	_, K := utils.SECP256k_Gen1G1KeyPair()
	v, V, _ := utils.BN254_GenG1KeyPair()

	r, R, _ := utils.BN254_GenG1KeyPair()

	S_Sender := ecpdksap_v2.SenderComputesSharedSecret(&r, &V, &K)

	S_Recipient := ecpdksap_v2.RecipientComputesSharedSecret(&v, &R, &K)
	
	if S_Sender != S_Recipient {
		t.Fatalf(`ERR: sender and recipient calculated different secret !!!`)
	}

	b_Sender := ecpdksap_v2.Compute_b(&S_Sender)
	b_Recipient := ecpdksap_v2.Compute_b(&S_Recipient)

	if b_Sender.String() != b_Recipient.String() {
		t.Fatalf(`ERR: sender and recipient calculated different 'b' !!!`)
	}
}
