package bn254_bench

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bls12-381"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	"ecpdksap-go/utils"
)

func Run(b *testing.B, sampleSize int, nRepetitions int, justViewTags bool, randomSeed int) map[string]time.Duration {

	fmt.Println("Running `bls12-381` Benchmark    ::: sampleSize:", sampleSize, "nRepetitions:", nRepetitions)
	fmt.Println()

	rndGen := rand.New(rand.NewSource(int64(randomSeed)))

	durations := map[string]time.Duration{}

	for iReps := 0; iReps < nRepetitions; iReps++ {

		g1, _, _, g2Aff := EC.Generators()

		//common for versions: V0, V1, V2
		_, v_asBigInt, V, _ := _EC_GenerateG1KeyPair(rndGen)
		v_asBigIntPtr := &v_asBigInt

		var r_asBigInt big.Int

		var P_v0 EC.GT

		//random data generation: Rj
		var Rs []EC.G1Jac
		var Rs_Ptr []*EC.G1Jac
		var RsAff_asArr [][]EC.G1Affine

		var rs []big.Int

		var combinedMeta []*_CombinedMeta

		for j := 0; j < sampleSize; j++ {

			_, rj_asBigInt, Rj, Rj_asAff := _EC_GenerateG1KeyPair(rndGen)

			Rs = append(Rs, Rj)
			RsAff_asArr = append(RsAff_asArr, []EC.G1Affine{Rj_asAff})

			tmp := new(EC.G1Jac)
			tmp.FromAffine(&Rj_asAff)
			Rs_Ptr = append(Rs_Ptr, tmp)

			//note: store the last priv. key for R
			r_asBigInt = rj_asBigInt
			rs = append(rs, r_asBigInt)

			cm := new(_CombinedMeta)
			cm.Rj = new(EC.G1Jac)
			cm.Rj.FromAffine(&Rj_asAff)
			cm.Rj_asAffArr = []EC.G1Affine{Rj_asAff}
			cm.ViewTagTwoBytes = uint16(rand.Uint32() % 65536)
			cm.ViewTagSingleByte = uint8(rand.Uint32() % 256)

			combinedMeta = append(combinedMeta, cm)
		}

		var rV EC.G1Jac
		rV.ScalarMultiplication(&V, &r_asBigInt)

		var rV_asAff EC.G1Affine
		rV_asAff.FromJacobian(&rV)

		//protocol V0 -------------------------------------

		_, _, _, K2_EC_asAff := _EC_GenerateG2KeyPair(rndGen)
		K2_EC_asAffArr := []EC.G2Affine{K2_EC_asAff}

		var vR EC.G1Jac
		var vR_asAff EC.G1Affine

		hasher := sha256.New()

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
		b.ResetTimer()

		for _, cm := range combinedMeta {

			hasher.Reset()

			compressed := (vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))).Bytes()

			if hasher.Sum(compressed[:])[0] != cm.ViewTagSingleByte {
				continue
			}

			pairingResult, _ := EC.Pair(cm.Rj_asAffArr, K2_EC_asAffArr)

			P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
		}

		durations["v0.v0-1byte"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-2bytes
		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagTwoBytes {
				continue
			}

			pairingResult, _ := EC.Pair(cm.Rj_asAffArr, K2_EC_asAffArr)

			P_v0.CyclotomicExp(pairingResult, v_asBigIntPtr)
		}

		durations["v0.v0-2bytes"] += b.Elapsed()

		//protocol: V0 and viewTag: V1-1byte
		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointXCoordToViewTagByte1(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagSingleByte {
				continue
			}

			pairingResult, _ := EC.Pair(cm.Rj_asAffArr, K2_EC_asAffArr)

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

		//protocol V2 --------------------

		_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
		var K_SECP256k1_Jac SECP256K1.G1Jac
		K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

		g2Aff_asArray := []EC.G2Affine{g2Aff}

		var Pv2_asJac SECP256K1.G1Jac

		b_asBigInt := new(big.Int)

		K_SECP256k1_JacPtr := &K_SECP256k1_Jac

		K_SECP256k1_AffPtr := new(SECP256K1.G1Affine)
		K_SECP256k1_AffPtr.FromJacobian(K_SECP256k1_JacPtr)

		//protocol: V2 and viewTag: none
		if !justViewTags {

			b.ResetTimer()

			for _, cm := range combinedMeta {

				S, _ := EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))}, g2Aff_asArray)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}

			durations["v2.none"] += b.Elapsed()
		}

		//protocol: V2 and viewTag: v0-1byte
		b.ResetTimer()

		for _, cm := range combinedMeta {

			hasher.Reset()

			compressed := (vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))).Bytes()

			if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

				S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

				Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
			}
		}

		durations["v2.v0-1byte"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-2bytes
		b.ResetTimer()

		for _, cm := range combinedMeta {

			if _EC_G1AffPointToViewTagByte2(vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))) != cm.ViewTagTwoBytes {
				continue
			}

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(b_asBigInt))
		}

		durations["v2.v0-2bytes"] += b.Elapsed()

		//protocol: V2 and viewTag: v1-1byte
		b.ResetTimer()

		for _, cm := range combinedMeta {

			if vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr)).Bytes()[0] != cm.ViewTagSingleByte {
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

	return durations
}

func _EC_GenerateG1KeyPair(r *rand.Rand) (privKey EC_fr.Element, privKey_asBigInt big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	randBigInt := big.NewInt(r.Int63())
	randBigInt.Mul(randBigInt, randBigInt).Mul(randBigInt, randBigInt)
	privKey.SetBigInt(randBigInt)

	privKey.BigInt(&privKey_asBigInt)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigInt)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _EC_GenerateG2KeyPair(r *rand.Rand) (privKey EC_fr.Element, privKey_asBigInt big.Int, pubKey EC.G2Jac, pubKeyAff EC.G2Affine) {
	_, g2, _, _ := EC.Generators()

	randBigInt := big.NewInt(r.Int63())
	randBigInt.Mul(randBigInt, randBigInt).Mul(randBigInt, randBigInt)
	privKey.SetBigInt(randBigInt)

	privKey.BigInt(&privKey_asBigInt)
	pubKey.ScalarMultiplication(&g2, &privKey_asBigInt)
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
}
