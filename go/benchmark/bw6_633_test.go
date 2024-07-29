package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bw6-633"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	"ecpdksap-go/utils"
)

func Benchmark_BW6_633(b *testing.B) {

	// _Benchmark_BW6_633(b, 5000, 10)
	// _Benchmark_BW6_633(b, 20_000, 10)
	// _Benchmark_BW6_633(b, 40_000, 10)
	_Benchmark_BW6_633(b, 80_000, 10)
	// _Benchmark_BW6_633(b, 100_000, 10)
}

func _Benchmark_BW6_633(b *testing.B, sampleSize int, nRepetitions int) {

	fmt.Println("Benchmark_BW6_633 ::: sampleSize:", sampleSize, "nRepetitions:", nRepetitions)
	fmt.Println()

	durations := map[string]time.Duration{}

	for i := 0; i < nRepetitions; i++ {

		g1, _, _, g2Aff := EC.Generators()

		//common for versions: V0, V1, V2
		_, v_asBigInt, V, _ := _BW6_633_GenerateG1KeyPair()

		var r_asBigInt big.Int

		var P_v0 EC.GT

		//random data generation: Rj
		var Rs []EC.G1Jac
		var RsAff_asArr [][]EC.G1Affine
		for j := 0; j < sampleSize; j++ {

			_, rj_asBigInt, Rj, Rj_asAff := _BW6_633_GenerateG1KeyPair()

			Rs = append(Rs, Rj)
			RsAff_asArr = append(RsAff_asArr, []EC.G1Affine{Rj_asAff})

			//note: store the last priv. key for R
			r_asBigInt = rj_asBigInt
		}

		//random data generation: view tags - 1 and 2 bytes
		var viewTags []uint16
		var viewTagsSingleByte []uint8

		//note: last view tag element will be overwritten by each tag's version
		for i := 0; i < sampleSize; i++ {
			tmp := rand.Uint32()
			viewTags = append(viewTags, uint16(tmp))
			viewTagsSingleByte = append(viewTagsSingleByte, uint8(tmp))
		}

		var rV EC.G1Jac
		rV.ScalarMultiplication(&V, &r_asBigInt)

		var rV_asAff EC.G1Affine
		rV_asAff.FromJacobian(&rV)

		//protocol V0 -------------------------------------

		_, _, _, K2_BW6_633_asAff := _BW6_633_GenerateG2KeyPair()
		K2_BW6_633_asAffArr := []EC.G2Affine{K2_BW6_633_asAff}

		var vR EC.G1Jac
		var vR_asAff EC.G1Affine

		//protocol: V0 and viewTag: none

		b.ResetTimer()

		for _, Rsi_asArray := range RsAff_asArr {

			pairingResult, _ := EC.Pair(Rsi_asArray, K2_BW6_633_asAffArr)

			P_v0.CyclotomicExp(pairingResult, &v_asBigInt)
		}

		durations["v0.none"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for i, Rsi_asArray := range RsAff_asArr {

			if _BW6_633_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rs[i], &v_asBigInt))) != viewTagsSingleByte[i] {
				continue
			}

			pairingResult, _ := EC.Pair(Rsi_asArray, K2_BW6_633_asAffArr)

			P_v0.CyclotomicExp(pairingResult, &v_asBigInt)
		}

		durations["v0.v0-1byte"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-2bytes
		viewTags[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for i, Rsi_asArray := range RsAff_asArr {

			if _BW6_633_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rs[i], &v_asBigInt))) != viewTags[i] {
				continue
			}

			pairingResult, _ := EC.Pair(Rsi_asArray, K2_BW6_633_asAffArr)

			P_v0.CyclotomicExp(pairingResult, &v_asBigInt)
		}

		durations["v0.v0-2bytes"] += b.Elapsed()

		//protocol: V0 and viewTag: V1-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for i, RsiAff_asArray := range RsAff_asArr {

			if _BW6_633_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rs[i], &v_asBigInt))) != viewTagsSingleByte[i] {
				continue
			}

			pairingResult, _ := EC.Pair(RsiAff_asArray, K2_BW6_633_asAffArr)

			P_v0.CyclotomicExp(pairingResult, &v_asBigInt)
		}

		durations["v0.v1-1byte"] += b.Elapsed()

		//protocol: V1 -------------------

		// var P_v1 EC.GT
		var tmp EC.G1Jac
		var tmpAff EC.G1Affine
		K_asArray := K2_BW6_633_asAffArr

		//protocol: V1 and viewTag: none
		b.ResetTimer()

		for _, Rsi_asJac := range Rs {

			vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt))

			hash_asBigInt := _BW6_633_HashG1AffPoint(&vR_asAff)

			tmp.ScalarMultiplication(&g1, hash_asBigInt)

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		}

		durations["v1.none"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for i, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(&vR)) != viewTagsSingleByte[i] {
				continue
			}

			hash_asBigInt := _BW6_633_HashG1AffPoint(&vR_asAff)

			tmp.ScalarMultiplication(&g1, hash_asBigInt)

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		}

		durations["v1.v0-1byte"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-2bytes
		viewTags[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for i, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(&vR)) != viewTags[i] {
				continue
			}

			hash_asBigInt := _BW6_633_HashG1AffPoint(&vR_asAff)

			tmp.ScalarMultiplication(&g1, hash_asBigInt)

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		}

		durations["v1.v0-2bytes"] += b.Elapsed()

		//protocol: V1 and viewTag: V1-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for i, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointXCoordToViewTagByte1(&vR_asAff) != viewTagsSingleByte[i] {
				continue
			}

			hash_asBigInt := _BW6_633_HashG1AffPoint(&vR_asAff)

			tmp.ScalarMultiplication(&g1, hash_asBigInt)

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		}

		durations["v1.v1-1byte"] += b.Elapsed()

		//protocol V2 --------------------

		_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
		var K_SECP256k1_Jac SECP256K1.G1Jac
		K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

		g2Aff_asArray := []EC.G2Affine{g2Aff}

		var Pv2_asJac SECP256K1.G1Jac

		b_asBigInt := new(big.Int)

		//protocol: V2 and viewTag: none
		b.ResetTimer()

		for _, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, g2Aff_asArray)

			//compute `b`
			S.B0.A0.BigInt(b_asBigInt)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, b_asBigInt)
		}

		durations["v2.none"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointToViewTagByte1(&vR_asAff) != viewTagsSingleByte[i] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, g2Aff_asArray)

			//compute `b`
			S.B0.A0.BigInt(b_asBigInt)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, b_asBigInt)
		}

		durations["v2.v0-1byte"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-2bytes
		viewTags[len(viewTags)-1] = _BW6_633_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointToViewTagByte2(&vR_asAff) != viewTags[i] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, g2Aff_asArray)

			//compute `b`
			S.B0.A0.BigInt(b_asBigInt)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, b_asBigInt)
		}

		durations["v2.v0-2bytes"] += b.Elapsed()

		//protocol: V2 and viewTag: v1-1byte
		viewTagsSingleByte[len(viewTags)-1] = _BW6_633_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs {

			vR.ScalarMultiplication(&Rsi_asJac, &v_asBigInt)

			if _BW6_633_G1AffPointXCoordToViewTagByte1(&vR_asAff) != viewTagsSingleByte[i] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, g2Aff_asArray)

			//compute `b`
			S.B0.A0.BigInt(b_asBigInt)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, b_asBigInt)
		}

		durations["v2.v1-1byte"] += b.Elapsed()
	}

	protocolVersions := []string{
		"v0.none", "v0.v0-1byte", "v0.v0-2bytes", "v0.v1-1byte",
		"v1.none", "v1.v0-1byte", "v1.v0-2bytes", "v1.v1-1byte",
		"v2.none", "v2.v0-1byte", "v2.v0-2bytes", "v2.v1-1byte",
	}

	for _, pVersion := range protocolVersions {
		fmt.Println("version:", pVersion, "duration:", durations[pVersion]/time.Duration(nRepetitions))
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
}

func _BW6_633_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _BW6_633_GenerateG2KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G2Jac, pubKeyAff EC.G2Affine) {
	_, g2, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g2, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _BW6_633_G1AffPointToViewTag(pt *EC.G1Affine, len uint) (viewTag string) {

	return _BW6_633_HashG1AffPoint(pt).Text(16)[:2*len]
}

func _BW6_633_G1AffPointToViewTagByte1(pt *EC.G1Affine) uint8 {
	hasher := sha256.New()
	compressed := pt.Bytes()
	hash := hasher.Sum(compressed[:])
	return hash[0]
}

func _BW6_633_G1AffPointToViewTagByte2(pt *EC.G1Affine) uint16 {
	hasher := sha256.New()
	hash := hasher.Sum(pt.Marshal())
	return binary.BigEndian.Uint16(hash[0:2])
}

func _BW6_633_G1AffPointXCoordToViewTag(pt *EC.G1Affine, len uint) (viewTag string) {

	return pt.X.Text(16)[:2*len]
}

func _BW6_633_G1AffPointXCoordToViewTagByte1(pt *EC.G1Affine) uint8 {
	return pt.X.Bytes()[0]
}

func _BW6_633_HashG1AffPoint(pt *EC.G1Affine) *big.Int {
	hasher := sha256.New()
	tmp := pt.X.Bytes()
	hasher.Write(tmp[:])
	tmp = pt.Y.Bytes()
	hasher.Write(tmp[:])
	hash_asBytes := hasher.Sum(nil)

	var hash EC_fr.Element

	hash_asBigInt := new(big.Int)

	hash.SetBytes(hash_asBytes)
	hash.BigInt(hash_asBigInt)

	return hash_asBigInt
}
