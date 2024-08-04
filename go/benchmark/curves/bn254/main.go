package bn254_bench

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"ecpdksap-go/utils"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"
)

func Run(b *testing.B, sampleSize int, nRepetitions int, justViewTags bool) {

	fmt.Println("Running `bn254` Benchmark ::: sampleSize:", sampleSize, "nRepetitions:", nRepetitions)
	fmt.Println()

	durations := map[string]time.Duration{}

	for iReps := 0; iReps < nRepetitions; iReps++ {

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

		var a, bx fp.Element

		var combinedMeta []*_CombinedMeta

		for j := 0; j < sampleSize; j++ {

			_, rj_asBigInt, Rj, Rj_asAff := _EC_GenerateG1KeyPair()

			Rs = append(Rs, Rj)
			RsAff_asArr = append(RsAff_asArr, []EC.G1Affine{Rj_asAff})

			tmp := new(EC.G1Jac)
			tmp.FromAffine(&Rj_asAff)
			Rs_Ptr = append(Rs_Ptr, tmp)

			//note: store the last priv. key for R
			r_asBigInt = rj_asBigInt
			rs = append(rs, r_asBigInt)

			//note: all tag elements will be overwritten by each protocol & tag's version


			cm := new(_CombinedMeta)
			cm.Rj = new(EC.G1Jac)
			cm.Rj.FromAffine(&Rj_asAff)
			cm.Rj_asAffArr = []EC.G1Affine{Rj_asAff}
			cm.ViewTagTwoBytes = uint16(rand.Uint32() % 65536)
			cm.ViewTagSingleByte = uint8(rand.Uint32() % 256)
			cm.ViewTagSecondByte = uint8(rand.Uint32() % 256)

			combinedMeta = append(combinedMeta, cm)
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

		vR.ScalarMultiplication(&Rs[len(Rs)-1], &v_asBigInt)

		fmt.Println("----->", rV.Equal(&vR))
		fmt.Println("MyFromJacobian2", rV.MyFromJacobian3() == vR.MyFromJacobian3())
		fmt.Println("MyFromJacobian2", rV.MyFromJacobian3(), vR.MyFromJacobian3())
		// pZSquare.Square(&p.Z)
		// aZSquare.Square(&q.Z)

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


		hasher := sha1.New()

		// precomputedLines := [][2][66]EC.LineEvaluationAff {EC.PrecomputeLines(K2_EC_asAffArr[0])}
		// precomputedLines := EC.PrecomputeLines(K2_EC_asAffArr[0])
		precomputedLines_asArr := [][2][66]EC.LineEvaluationAff {EC.PrecomputeLines(K2_EC_asAffArr[0])}

		// precomputedLinesPtr := &precomputedLines

		e := P_v0.MyE(v_asBigIntPtr)

		var Rjs []*EC.G1Jac

		for _, cm := range combinedMeta {
			Rjs = append(Rjs, cm.Rj)
		}

		b.ResetTimer()

		vRs := vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)

		for i, cm := range combinedMeta {

			hasher.Reset()

			preimage := vRs[i].MyFromJacobian3()
			compressed := preimage.Bytes()

			if hasher.Sum(compressed[:])[0] != cm.ViewTagSingleByte {
				continue
			}

			pairingResult, _ := EC.PairFixedQ(cm.Rj_asAffArr, precomputedLines_asArr)

			P_v0.MyCyclotomicExp(pairingResult, e)
		}

		durations["v0.v0-1byte"] += b.Elapsed()


		//protocol: V0 and viewTag: V0-2bytes

		b.ResetTimer()

		vRs = vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)

		for i, cm := range combinedMeta {

			hasher.Reset()

			preimage := vRs[i].MyFromJacobian3()
			compressed := preimage.Bytes()

			vTag := hasher.Sum(compressed[:])

			if vTag[0] != cm.ViewTagSingleByte || vTag[1] != cm.ViewTagSecondByte {
				continue
			}

			pairingResult, _ := EC.PairFixedQ(cm.Rj_asAffArr, precomputedLines_asArr)

			P_v0.MyCyclotomicExp(pairingResult, e)
		}

		durations["v0.v0-2bytes"] += b.Elapsed()

		// //protocol: V0 and viewTag: V1-1byte

		b.ResetTimer()

		vRs = vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)
		
		for i, vrs := range vRs {

			if vR_asAff.X.Mul(&vrs.X, bx.Square(a.Inverse(&vrs.Z))).Bytes()[0] == combinedMeta[i].ViewTagSingleByte {

				pairingResult, _ := EC.PairFixedQ(combinedMeta[i].Rj_asAffArr, precomputedLines_asArr)

				P_v0.MyCyclotomicExp(pairingResult, e)
			}
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

				tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(vR_asAff.FromJacobian(vR.ScalarMultiplication(&Rsi_asJac, v_asBigIntPtr))))

				EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
			}

			durations["v1.none"] += b.Elapsed()
		}
		//protocol: V1 and viewTag: V0-1byte

		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagSingleByte {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v0-1byte"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-2bytes

		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagTwoBytes {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v0-2bytes"] += b.Elapsed()

		//protocol: V1 and viewTag: V1-1byte

		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagSingleByte {
				continue
			}

			EC.Pair([]EC.G1Affine{*tmpAff.FromJacobian(tmp.ScalarMultiplication(&g1, _EC_HashG1AffPoint(&vR_asAff)))}, K_asArray)
		}

		durations["v1.v1-1byte"] += b.Elapsed()

		// //protocol V2 --------------------

		_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
		var K_SECP256k1_Jac SECP256K1.G1Jac
		K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

		g2Aff_asArray := []EC.G2Affine{g2Aff}

		var Pv2_asJac SECP256K1.G1Jac

		b_asBigInt := new(big.Int)

		K_SECP256k1_JacPtr := &K_SECP256k1_Jac

		K_SECP256k1_AffPtr := new(SECP256K1.G1Affine)
		K_SECP256k1_AffPtr.FromJacobian(K_SECP256k1_JacPtr)
		// var bs []SECP256K1_fr.Element

		// //protocol: V2 and viewTag: none
		// if !justViewTags {

		// 	// bs = []SECP256K1_fr.Element{}

		// 	b.ResetTimer()

		// 	for _, cm := range combinedMeta {

		// 		S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))}, g2Aff_asArray)

		// 		// bs = append(bs, SECP256K1_fr.Element(S.C0.B0.A0))

		// 		Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
		// 	}

		// 	// SECP256K1.BatchScalarMultiplicationG1(K_SECP256k1_AffPtr, bs)
		// 	durations["v2.none"] += b.Elapsed()
		// }

		// //protocol: V2 and viewTag: v0-1byte

		precomputedLines_asArr = [][2][66]EC.LineEvaluationAff {EC.PrecomputeLines(g2Aff_asArray[0])}


		b.ResetTimer()

		vRs = vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)
		
		for i, cm := range combinedMeta {

			hasher.Reset()

			compressed := vR_asAff.X.Mul(&vRs[i].X, bx.Square(a.Inverse(&vRs[i].Z))).Bytes()

			if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

				vR_asAff.Y.Mul(&vRs[i].Y, &bx).Mul(&vR_asAff.Y, &a)

				S, _ := EC.PairFixedQ([]EC.G1Affine{vR_asAff}, precomputedLines_asArr)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}
		}

		durations["v2.v0-1byte"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-2bytes
		b.ResetTimer()

		vRs = vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)

		for i, cm := range combinedMeta {

			hasher.Reset()

			compressed := vR_asAff.X.Mul(&vRs[i].X, bx.Square(a.Inverse(&vRs[i].Z))).Bytes()

			vTag := hasher.Sum(compressed[:])

			if vTag[0] == cm.ViewTagSingleByte && vTag[1] == cm.ViewTagSecondByte {

				vR_asAff.Y.Mul(&vRs[i].Y, &bx).Mul(&vR_asAff.Y, &a)

				S, _ := EC.PairFixedQ([]EC.G1Affine{vR_asAff}, precomputedLines_asArr)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}
		}

		durations["v2.v0-2bytes"] += b.Elapsed()

		// //protocol: V2 and viewTag: v1-1byte
		b.ResetTimer()

		vRs = vR.BatchScalarMultiplicationUsingFixedS(&Rjs, v_asBigIntPtr)
		
		for i, vrs := range vRs {

			if vR_asAff.X.Mul(&vrs.X, bx.Square(a.Inverse(&vrs.Z))).Bytes()[0] == combinedMeta[i].ViewTagSingleByte {

				vR_asAff.Y.Mul(&vrs.Y, &bx).Mul(&vR_asAff.Y, &a)

				S, _ := EC.PairFixedQ([]EC.G1Affine{vR_asAff}, precomputedLines_asArr)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}
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

func _EC_G1AffPointToViewTagByte1(pt *EC.G1Affine) uint8 {
	hasher := sha256.New()
	compressed := pt.Bytes()
	return hasher.Sum(compressed[:])[0]
}

func _EC_G1AffPointToViewTagByte2(pt *EC.G1Affine) uint16 {
	hasher := sha256.New()
	compressed := pt.Bytes()
	return binary.BigEndian.Uint16(hasher.Sum(compressed[:])[0:2])
}

func _EC_G1AffPointXCoordToViewTagByte1(pt *EC.G1Affine) uint8 {
	return pt.X.Bytes()[0]
}

func _EC_HashG1AffPoint(pt *EC.G1Affine) *big.Int {
	hasher := sha256.New()
	compressed := pt.Bytes()

	var hash EC_fr.Element

	return hash.SetBytes(hasher.Sum(compressed[:])).BigInt(new(big.Int))
}

type _CombinedMeta struct {
	Rj                *EC.G1Jac
	Rj_asAffArr       []EC.G1Affine
	ViewTagTwoBytes   uint16
	ViewTagSingleByte uint8
	ViewTagSecondByte   uint8

}
