package main

import (
	"ecpdksap-bn254/gen_example"
	"ecpdksap-bn254/recipient"
	"ecpdksap-bn254/sender"
	"ecpdksap-bn254/utils"
	ecpdksap_v2 "ecpdksap-bn254/versions/v2"
	"fmt"
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
	SECP256K1_fr "github.com/consensys/gnark-crypto/ecc/secp256k1/fr"
)

func DBG() {
	privK, pubK := utils.GenSECP256k1G1KeyPair()

	addr := ecpdksap_v2.ComputeEthAddress(&pubK)

	fmt.Println("privK:", "0x"+privK.Text(16), "\n")
	fmt.Println("pubK:", pubK.X.Text(16)+"."+pubK.Y.Text(16))
	fmt.Println("addr:", addr)
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("\nERR: All subcommands (first arg.) receive only one param (second arg.).")
		return
	}

	subcmd := os.Args[1]
	arg := os.Args[2]

	switch subcmd {

	case "send":
		sender.Send(arg)

	case "receive-scan":
		_, addrs, privKeys := recipient.Scan(arg)

		fmt.Println(addrs[0], "\n", privKeys[0])

	case "receive-scan-using-vtag":
		rP := recipient.ScanUsingViewTag(arg)
		fmt.Println(rP)

	case "gen-example":
		gen_example.GenerateExample("v2")

	case "dbg":
		DBG()


	default:
		fmt.Printf("\nERR: only: `send` | `receive-scan` | `receive-scan-using-vtag` subcommands allowed.\n\n")
		return
	}
}



func DBG2() {
	one := new(big.Int)
	one.SetString("1", 10)

	var p big.Int
	p.SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10) 

	var k_asElement SECP256K1_fr.Element
	k_asElement.SetRandom()
	k := new(big.Int)
	k_asElement.BigInt(k)
	k.Mod(k, &p)


	// k.SetString("1123", 10)

	var v_asElement fr.Element
	v_asElement.SetRandom()
	v := new(big.Int)
	v_asElement.BigInt(v)
	v.Mod(v, &p)

	// v.SetString("3333", 10)

	var r_asElement fr.Element
	r_asElement.SetRandom()
	r := new(big.Int)
	r_asElement.BigInt(r)
	r.Mod(r, &p)

	//---------------------------------------------------

	_, G1_SECP256K1 := SECP256K1.Generators()

	var K SECP256K1.G1Affine
	K.ScalarMultiplicationBase(k)

	V, _ := utils.CalcG1PubKey(v_asElement)

	R, _ := utils.CalcG1PubKey(r_asElement)

	var rV bn254.G1Affine
	rV.ScalarMultiplication(&V, r)

	var vR bn254.G1Affine
	vR.ScalarMultiplication(&R, v)

	// Compute pairing
	_, _, _, G2_BN254 := bn254.Generators()


	P_for_b, _ := bn254.Pair([]bn254.G1Affine{rV}, []bn254.G2Affine{G2_BN254})
	P_for_b_recipient, _ := bn254.Pair([]bn254.G1Affine{vR}, []bn254.G2Affine{G2_BN254})

	fmt.Println("P_for_b is matches one both sides:", P_for_b == P_for_b_recipient)

	var res bn254.E2
	var b big.Int

	res.Add(&P_for_b.C0.B0, &P_for_b.C0.B1)
	res.Add(&res, &P_for_b.C1.B0)
	res.Add(&res, &P_for_b.C1.B1)

	res.A0.BigInt(&b)
	b.Add(&b, res.A1.BigInt(new(big.Int)))

	// b.Mod(&b, &p)

	// b.SetString("100000000000", 10)

	var b_asElement SECP256K1_fr.Element
	b_asElement.SetRandom()

	var kb_asElement SECP256K1_fr.Element

	kb_asElement.Mul(&k_asElement, &b_asElement)

	var kb_asBigInt big.Int
	kb_asElement.BigInt(&kb_asBigInt)

	var b_asBigInt big.Int
	// b_asBigInt
	b_asElement.BigInt(&b_asBigInt)

	var bK SECP256K1.G1Affine
	bK.ScalarMultiplication(&K, &b_asBigInt)

	var bkG1 SECP256K1.G1Affine
	var bk big.Int
	bk.Mul(&b, k) // bk > p ----> M ; bk % p ----> M2
	// bk.Mod(&bk, &p)

	// bkG1.ScalarMultiplication(&G1_SECP256K1, &b)
	// bkG1.ScalarMultiplication(&bkG1, k)
	bkG1.ScalarMultiplication(&G1_SECP256K1, &kb_asBigInt)

	fmt.Println("P matches on both sides:", bK == bkG1)

	// fmt.Println("bK:", bK)
	// fmt.Println("bkG1:", bkG1)

	var P_using_base SECP256K1.G1Affine
	P_using_base.ScalarMultiplicationBase(&kb_asBigInt)

	fmt.Println("P matches with bk*G1 using base mul.:", bK == P_using_base)

	// bk.Mod(&bk, &p)
	fmt.Println("")
	fmt.Println("bk:", kb_asBigInt.Text(16))
	fmt.Println("P.X", bK.X.Text(16))
	fmt.Println("P.Y", bK.Y.Text(16))
	fmt.Println("Is on curve:", bK.IsOnCurve())
}



