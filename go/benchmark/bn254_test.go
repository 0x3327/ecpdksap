package main

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	BN254 "github.com/consensys/gnark-crypto/ecc/bn254"
	BN254_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

	ecpdksap_v2 "ecpdksap-go/versions/v2"

	"ecpdksap-go/utils"
)

func Benchmark_BN254(b *testing.B) {

	_Benchmark_BN254(b, 1000, 10)

	// Benchmark_BN254(100_000, 10)
}


func _Benchmark_BN254(b *testing.B, sampleSize int, nRepetitions int) {

	fmt.Println("Benchmark_BN254 ::: sampleSize:", sampleSize, "nRepetitions:", nRepetitions)
	fmt.Println()

    durations := map[string]time.Duration{}

	for i := 0; i < nRepetitions; i++ {

		g1, _, _, g2Aff := BN254.Generators()

		//common for versions: V0, V1, V2
		v, V_aff, _ := utils.BN254_GenG1KeyPair()
		var V BN254.G1Jac
		V.FromAffine(&V_aff)
		v_asBigInt := new(big.Int)
		v.BigInt(v_asBigInt)

		var r_asBigInt big.Int

		var P_v0 BN254.GT

		//random data generation
		var Rs [] BN254.G1Jac
		var RsAff_asArr [][] BN254.G1Affine
		for i := 0; i < sampleSize; i++ {
			var ri BN254_fr.Element
			ri.SetRandom()
			var ri_asBigInt big.Int
			ri.BigInt(&ri_asBigInt)
			var Ri BN254.G1Jac
			Ri.ScalarMultiplication(&g1, &ri_asBigInt)
			Rs = append(Rs, Ri)
			
			ri.BigInt(&r_asBigInt)

			var Ri_asAff BN254.G1Affine
			Ri_asAff.FromJacobian(&Ri)
			RsAff_asArr = append(RsAff_asArr, []BN254.G1Affine{Ri_asAff})
		}

		var rV BN254.G1Jac
		rV.ScalarMultiplication(&V, &r_asBigInt)

		var viewTags []string

		for i := 0; i < sampleSize; i++ {
			_, pt_rand, _ := utils.BN254_GenG1KeyPair()
			viewTags = append(viewTags, utils.BN254_G1PointToViewTag(&pt_rand, 2))
		}

		//protocol V0 -------------------------------------

		_, K_G2BN254, _ := utils.BN254_GenG2KeyPair()
		K_G2BN254_asArray := []BN254.G2Affine{K_G2BN254}

		var vR_asJac BN254.G1Jac

		//protocol: V0 and viewTag: none

		b.ResetTimer()

		for _, Rsi_asArray := range RsAff_asArr { 
					
			pairingResult, _ := BN254.Pair(Rsi_asArray, K_G2BN254_asArray)

			P_v0.CyclotomicExp(pairingResult, v_asBigInt)
		}

		durations["v0.none"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 1)

		b.ResetTimer()

		for i, Rsi_asArray := range RsAff_asArr { 

			if utils.BN254_G1JacPointToViewTag(vR_asJac.ScalarMultiplication(&Rs[i], v_asBigInt), 1) != viewTags[i][:2] {
				continue
			}
					
			pairingResult, _ := BN254.Pair(Rsi_asArray, K_G2BN254_asArray)

			P_v0.CyclotomicExp(pairingResult, v_asBigInt)
		}

		durations["v0.v0-1byte"] += b.Elapsed()

		//protocol: V0 and viewTag: V0-2bytes
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 2)

		b.ResetTimer()

		for i, Rsi_asArray := range RsAff_asArr { 
	
			if utils.BN254_G1JacPointToViewTag(vR_asJac.ScalarMultiplication(&Rs[i], v_asBigInt), 2) != viewTags[i] {
				continue
			}
					
			pairingResult, _ := BN254.Pair(Rsi_asArray, K_G2BN254_asArray)
	
			P_v0.CyclotomicExp(pairingResult, v_asBigInt)
		}
	
		durations["v0.v0-2bytes"] += b.Elapsed()

		//protocol: V0 and viewTag: V1-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointXCoordToViewTag(&rV, 1)

		b.ResetTimer()

		for i, Rsi_asArray := range RsAff_asArr { 
	
			if utils.BN254_G1JacPointXCoordToViewTag(vR_asJac.ScalarMultiplication(&Rs[i], v_asBigInt), 1) != viewTags[i][:2] {
				continue
			}
					
			pairingResult, _ := BN254.Pair(Rsi_asArray, K_G2BN254_asArray)
	
			P_v0.CyclotomicExp(pairingResult, v_asBigInt)
		}
	
		durations["v0.v1-1byte"] += b.Elapsed()

		//protocol: V1 -------------------

		// var P_v1 BN254.GT
		var hash BN254_fr.Element
		var hash_asBigInt big.Int
		var tmp BN254.G1Jac
		var tmpAff BN254.G1Affine
		K_asArray := []BN254.G2Affine{K_G2BN254}

		//protocol: V1 and viewTag: none
		b.ResetTimer()

		for _, Rsi_asJac := range Rs { 
				
			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)

			hash_asBytes := utils.BN254_HashG1JacPoint(&vR_asJac)
			hash.SetBytes(hash_asBytes)
			hash.BigInt(&hash_asBigInt)

			tmp.ScalarMultiplication(&g1, &hash_asBigInt)
			
			BN254.Pair([]BN254.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		} 

		durations["v1.none"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 1)

		b.ResetTimer()

		for i, Rsi_asJac := range Rs { 
			
			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)
			
			if utils.BN254_G1JacPointToViewTag(&vR_asJac, 1) != viewTags[i][:2] { continue }

			hash_asBytes := utils.BN254_HashG1JacPoint(&vR_asJac)
			hash.SetBytes(hash_asBytes)
			hash.BigInt(&hash_asBigInt)

			tmp.ScalarMultiplication(&g1, &hash_asBigInt)
			
			BN254.Pair([]BN254.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		} 

		durations["v1.v0-1byte"] += b.Elapsed()

		//protocol: V1 and viewTag: V0-2bytes
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 2)

		b.ResetTimer()

		for i, Rsi_asJac := range Rs { 
			
			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)
			
			if utils.BN254_G1JacPointToViewTag(&vR_asJac, 2) != viewTags[i] { continue }

			hash_asBytes := utils.BN254_HashG1JacPoint(&vR_asJac)
			hash.SetBytes(hash_asBytes)
			hash.BigInt(&hash_asBigInt)

			tmp.ScalarMultiplication(&g1, &hash_asBigInt)
			
			BN254.Pair([]BN254.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		} 

		durations["v1.v0-2bytes"] += b.Elapsed()

		//protocol: V1 and viewTag: V1-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointXCoordToViewTag(&rV, 1)
		
		b.ResetTimer()
		
		for i, Rsi_asJac := range Rs { 
			
			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)
			
			if utils.BN254_G1JacPointXCoordToViewTag(&vR_asJac, 1) != viewTags[i][:2] { continue }

			hash_asBytes := utils.BN254_HashG1JacPoint(&vR_asJac)
			hash.SetBytes(hash_asBytes)
			hash.BigInt(&hash_asBigInt)

			tmp.ScalarMultiplication(&g1, &hash_asBigInt)
			
			BN254.Pair([]BN254.G1Affine{*tmpAff.FromJacobian(&tmp)}, K_asArray)
		} 

		durations["v1.v1-1byte"] += b.Elapsed()

		//protocol V2 --------------------

		_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
		var K_SECP256k1_Jac SECP256K1.G1Jac
		K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

		var vR BN254.G1Affine
		g2Aff_asArray := []BN254.G2Affine{g2Aff}

		var Pv2_asJac SECP256K1.G1Jac
		//protocol: V2 and viewTag: none

		b.ResetTimer()

		for _, Rsi_asJac := range Rs { 

			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)

			S, _ := BN254.Pair([]BN254.G1Affine{*vR.FromJacobian(&vR_asJac)}, g2Aff_asArray)
			// b := ecpdksap_v2.Compute_b(&S)

			// Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, &b)

			b := ecpdksap_v2.Compute_b_asElement(&S)
			utils.SECP256k1_MulG1JacPointandElement(&K_SECP256k1_Jac, &b)
		}

		durations["v2.none"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 1)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs { 

			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)

			if utils.BN254_G1JacPointToViewTag(&vR_asJac, 1) != viewTags[i][:2] { continue }

			S, _ := BN254.Pair([]BN254.G1Affine{*vR.FromJacobian(&vR_asJac)}, g2Aff_asArray)
			b := ecpdksap_v2.Compute_b(&S)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, &b)
		}

		durations["v2.v0-1byte"] += b.Elapsed()

		//protocol: V2 and viewTag: v0-2bytes
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointToViewTag(&rV, 2)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs { 

			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)

			if utils.BN254_G1JacPointToViewTag(&vR_asJac, 2) != viewTags[i] { continue }

			S, _ := BN254.Pair([]BN254.G1Affine{*vR.FromJacobian(&vR_asJac)}, g2Aff_asArray)
			b := ecpdksap_v2.Compute_b(&S)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, &b)
		}

		durations["v2.v0-2bytes"] += b.Elapsed()

		//protocol: V2 and viewTag: v1-1byte
		viewTags[len(viewTags)-1] = utils.BN254_G1JacPointXCoordToViewTag(&rV, 1)

		b.ResetTimer()

		for _, Rsi_asJac := range Rs { 

			vR_asJac.ScalarMultiplication(&Rsi_asJac, v_asBigInt)

			if utils.BN254_G1JacPointXCoordToViewTag(&vR_asJac, 1) != viewTags[i][:2] { continue }

			S, _ := BN254.Pair([]BN254.G1Affine{*vR.FromJacobian(&vR_asJac)}, g2Aff_asArray)
			b := ecpdksap_v2.Compute_b(&S)

			Pv2_asJac.ScalarMultiplication(&K_SECP256k1_Jac, &b)
		}

		durations["v2.v1-1byte"] += b.Elapsed()
	}

	protocolVersions := []string {
		"v0.none", "v0.v0-1byte", "v0.v0-2bytes", "v0.v1-1byte",
		"v1.none", "v1.v0-1byte", "v1.v0-2bytes", "v1.v1-1byte",
		"v2.none", "v2.v0-1byte", "v2.v0-2bytes", "v2.v1-1byte",
	}

	for _, pVersion := range protocolVersions {
		fmt.Println("version:", pVersion, "duration:", durations[pVersion] / time.Duration(nRepetitions))
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
}

