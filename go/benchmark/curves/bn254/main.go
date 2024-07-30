package bn254_bench

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	"ecpdksap-go/utils"
)

func Run(b *testing.B, sampleSize int, nRepetitions int, justViewTags bool) {

	fmt.Println("Running `bn254` Benchmark ::: sampleSize:", sampleSize, "nRepetitions:", nRepetitions)
	fmt.Println()

	durations := map[string]time.Duration{}

	for i := 0; i < nRepetitions; i++ {

		g1, _, _, g2Aff := EC.Generators()

		//common for versions: V0, V1, V2
		_, v_asBigInt, V, _ := _EC_GenerateG1KeyPair()
		v_asBigIntPtr := &v_asBigInt

		var r_asBigInt big.Int

		var P_v0 EC.GT

		//random data generation: Rj
		var Rs []EC.G1Jac
		var Rs_Ptr []*EC.G1Jac
		var RsAff_asArr [][]EC.G1Affine

		var rs []big.Int

		var viewTagsTwoBytes []uint16
		var viewTagsSingleByte []uint8

		for j := 0; j < sampleSize; j++ {

			_, rj_asBigInt, Rj, Rj_asAff := _EC_GenerateG1KeyPair()

			Rs = append(Rs, Rj)
			RsAff_asArr = append(RsAff_asArr, []EC.G1Affine{Rj_asAff})

			tmp := new (EC.G1Jac)
			tmp.FromAffine(&Rj_asAff)
			Rs_Ptr = append(Rs_Ptr, tmp)

			//note: store the last priv. key for R
			r_asBigInt = rj_asBigInt
			rs = append(rs, r_asBigInt)

			//note: all tag elements will be overwritten by each protocol & tag's version
			
			viewTagsTwoBytes = append(viewTagsTwoBytes, uint16(rand.Uint32()%65536))
			viewTagsSingleByte = append(viewTagsSingleByte, uint8(rand.Uint32() % 256))
		}

		var rV EC.G1Jac
		rV.ScalarMultiplication(&V, &r_asBigInt)

		var rV_asAff EC.G1Affine
		rV_asAff.FromJacobian(&rV)

		//protocol V0 -------------------------------------

		_, _, _, K2_EC_asAff := _EC_GenerateG2KeyPair()
		K2_EC_asAffArr := []EC.G2Affine{K2_EC_asAff}

		var vR EC.G1Jac
		var vR_asAff EC.G1Affine

		//protocol: V0 and viewTag: none
		if !justViewTags {

			b.ResetTimer()

			for _, Rsi_asArray := range RsAff_asArr {

				pairingResult, _ := EC.Pair(Rsi_asArray, K2_EC_asAffArr)

				P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
			}

			durations["v0.none"] += b.Elapsed()
		}
		//protocol: V0 and viewTag: V0-1byte

		viewTagsSingleByte[len(viewTagsTwoBytes)-1] = _EC_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			pairingResult, _ := EC.Pair(RsAff_asArr[j], K2_EC_asAffArr)

			P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
		}

		durations["v0.v0-1byte"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-2bytes
		viewTagsTwoBytes[len(viewTagsTwoBytes)-1] = _EC_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsTwoBytes[j] {
				continue
			}

			pairingResult, _ := EC.Pair(RsAff_asArr[j], K2_EC_asAffArr)

			P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
		}

		durations["v0.v0-2bytes"] += b.Elapsed()

		//protocol: V0 and viewTag: V1-1byte
		viewTagsSingleByte[len(viewTagsSingleByte)-1] = _EC_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			pairingResult, _ := EC.Pair(RsAff_asArr[j], K2_EC_asAffArr)

			P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
		}

		durations["v0.v1-1byte"] += b.Elapsed()

		//protocol: V1 -------------------

		// var P_v1 EC.GT
		var tmp EC.G1Jac
		var tmpAff EC.G1Affine
		K_asArray := K2_EC_asAffArr

		//protocol: V1 and viewTag: none
		if !justViewTags {
			b.ResetTimer()

			for _, Rsi_asJac := range Rs {

				vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rsi_asJac, v_asBigIntPtr))

				hash_asBigInt := _EC_HashG1AffPoint(&vR_asAff)

				tmp.ScalarMultiplication(&g1, hash_asBigInt)

				EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
			}

			durations["v1.none"] += b.Elapsed()
		}
		//protocol: V1 and viewTag: V0-1byte
		viewTagsSingleByte[len(viewTagsSingleByte)-1] = _EC_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v0-1byte"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-2bytes
		viewTagsTwoBytes[len(viewTagsTwoBytes)-1] = _EC_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsTwoBytes[j] {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v0-2bytes"] += b.Elapsed()

		//protocol: V1 and viewTag: V1-1byte
		viewTagsSingleByte[len(viewTagsSingleByte)-1] = _EC_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v1-1byte"] += b.Elapsed()

		//protocol V2 --------------------

		_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
		var K_SECP256k1_Jac SECP256K1.G1Jac
		K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

		g2Aff_asArray := []EC.G2Affine{g2Aff}

		var Pv2_asJac SECP256K1.G1Jac

		b_asBigInt := new(big.Int)

		K_SECP256k1_JacPtr := &K_SECP256k1_Jac

		//protocol: V2 and viewTag: none
		if !justViewTags {

			b.ResetTimer()

			for _, Rsi_asJac := range Rs {

				S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rsi_asJac, v_asBigIntPtr))}, g2Aff_asArray)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}
			durations["v2.none"] += b.Elapsed()
		}

		//protocol: V2 and viewTag: v0-1byte
		viewTagsSingleByte[len(viewTagsSingleByte)-1] = _EC_G1AffPointToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
		}

		durations["v2.v0-1byte"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-2bytes
		viewTagsTwoBytes[len(viewTagsTwoBytes)-1] = _EC_G1AffPointToViewTagByte2(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsTwoBytes[j] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
		}

		durations["v2.v0-2bytes"] += b.Elapsed()

		//protocol: V2 and viewTag: v1-1byte
		viewTagsSingleByte[len(viewTagsSingleByte)-1] = _EC_G1AffPointXCoordToViewTagByte1(&rV_asAff)

		b.ResetTimer()

		for j, Rsj := range Rs_Ptr {

			if _EC_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(Rsj, v_asBigIntPtr))) != viewTagsSingleByte[j] {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
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

func _EC_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _EC_GenerateG2KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G2Jac, pubKeyAff EC.G2Affine) {
	_, g2, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g2, &privKey_asBigIng)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _EC_G1AffPointToViewTag(pt *EC.G1Affine, len uint) (viewTag string) {

	return _EC_HashG1AffPoint(pt).Text(16)[:2*len]
}

func _EC_G1AffPointToViewTagByte1(pt *EC.G1Affine) uint8 {
	hasher := sha256.New()
	compressed := pt.Bytes()
	hash := hasher.Sum(compressed[:])
	return hash[0]
}

func _EC_G1AffPointToViewTagByte2(pt *EC.G1Affine) uint16 {
	hasher := sha256.New()
	compressed := pt.Bytes()
	hash := hasher.Sum(compressed[:])
	return binary.BigEndian.Uint16(hash[0:2])
}

func _EC_G1AffPointXCoordToViewTag(pt *EC.G1Affine, len uint) (viewTag string) {

	return pt.X.Text(16)[:2*len]
}

func _EC_G1AffPointXCoordToViewTagByte1(pt *EC.G1Affine) uint8 {
	return pt.X.Bytes()[0]
}

func _EC_HashG1AffPoint(pt *EC.G1Affine) *big.Int {
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