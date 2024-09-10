package benchmark

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
<<<<<<< HEAD

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fp "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

	SECP256K1 "github.com/consensys/gnark-crypto/ecc/secp256k1"

=======
	"time"

	EC "github.com/consensys/gnark-crypto/ecc/bn254"
	EC_fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"

>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
	bls12_377 "ecpdksap-go/benchmark/curves/bls12-377"
	bls12_381 "ecpdksap-go/benchmark/curves/bls12-381"
	bls24_315 "ecpdksap-go/benchmark/curves/bls24-315"
	bn254 "ecpdksap-go/benchmark/curves/bn254"
	bw6_633 "ecpdksap-go/benchmark/curves/bw6-633"
	bw6_761 "ecpdksap-go/benchmark/curves/bw6-761"

<<<<<<< HEAD
	"ecpdksap-go/utils"
)

func Benchmark_BN254_V2_V0_1byte_Combined(b *testing.B, sampleSize int) {

	hasher := sha256.New()
	combinedMeta, v_asBigIntPtr, _, K_SECP256k1_JacPtr := _generateData(sampleSize)
	var vR EC.G1Jac
	var vR_asAff EC.G1Affine
	_, _, _, g2Aff := EC.Generators()
	g2Aff_asArray := []EC.G2Affine{g2Aff}
	var Pv2_asJac SECP256K1.G1Jac
	
	b.ResetTimer()
	for _, cm := range combinedMeta {

		hasher.Reset()

		compressed := (vR_asAff.FromJacobian(vR.ScalarMultiplication(cm.Rj, v_asBigIntPtr))).X.Bytes()

		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(new (big.Int)))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExternalCalls) :: Total time:", b.Elapsed())

	neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix := EC.PrecomputationForFixedScalarMultiplication(v_asBigIntPtr)
	var table [15]EC.G1Jac
	var a_El, b_El *EC_fp.Element
	var b_asBigInt big.Int

	b.ResetTimer()
	for _, cm := range combinedMeta {

		hasher.Reset()

		vR.FixedScalarMultiplication(cm.Rj, &table, neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix)

		a_El, b_El = vR_asAff.FromJacobianCoordX(&vR)

		compressed := (vR_asAff).X.Bytes()
		
		if hasher.Sum(compressed[:])[0] == cm.ViewTagSingleByte {

			vR_asAff.FromJacobianCoordY(a_El, b_El, &vR)

			S, _ := EC.Pair([]EC.G1Affine{vR_asAff}, g2Aff_asArray)

			Pv2_asJac.ScalarMultiplication(K_SECP256k1_JacPtr, S.C0.B0.A0.BigInt(&b_asBigInt))
		}
	}

	fmt.Println("(_Benchmark_BN254_V2_V0_1byte_ExpandedGnarkCrypto) :: Total time:", b.Elapsed())
}


func _generateData (sampleSize int) (combinedMeta  []*_CombinedMeta, v_asBigIntPtr *big.Int, V_Ptr EC.G1Jac, K_SECP256k1_JacPtr *SECP256K1.G1Jac) {
	_, v_asBigInt, V, _ := _EC_GenerateG1KeyPair()
	v_asBigIntPtr = &v_asBigInt
	V_Ptr = V

	var r_asBigInt big.Int

	//random data generation: Rj
	var Rs []EC.G1Jac
	var Rs_Ptr []*EC.G1Jac
	var RsAff_asArr [][]EC.G1Affine

	var rs []big.Int

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

		cm := new(_CombinedMeta)
		cm.Rj = new(EC.G1Jac)
		cm.Rj.FromAffine(&Rj_asAff)
		cm.Rj_asAffArr = []EC.G1Affine{Rj_asAff}
		cm.ViewTagTwoBytes = uint16(rand.Uint32() % 65536)
		cm.ViewTagSingleByte = uint8(rand.Uint32() % 256)

		combinedMeta = append(combinedMeta, cm)
	}


	_, K_SECP256k1 := utils.SECP256k_Gen1G1KeyPair()
	var K_SECP256k1_Jac SECP256K1.G1Jac
	K_SECP256k1_Jac.FromAffine(&K_SECP256k1)

	K_SECP256k1_JacPtr = &K_SECP256k1_Jac

	return
}


func _EC_GenerateG1KeyPair() (privKey EC_fr.Element, privKey_asBigIng big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey.SetRandom()
	privKey.BigInt(&privKey_asBigIng)
	pubKey.ScalarMultiplication(&g1, &privKey_asBigIng)
=======
	bn254_optimized "ecpdksap-go/benchmark/bn254"
	bn254_crk "ecpdksap-go/benchmark/bn254_constant_recipient_keys"
	bn254_v2_wea "ecpdksap-go/benchmark/bn254_v2_without_eth_addr"
)

func RunBench(kind string, rndSeed int) {

	b := new(testing.B)
	b.StartTimer()

	if kind == "only-bn254" {

		bn254_optimized.Run(b, 5_000, 10, rndSeed)
		bn254_optimized.Run(b, 10_000, 10, rndSeed)
		// bn254_optimized.Run(b, 20_000, 10, rndSeed)
		// bn254_optimized.Run(b, 40_000, 10, rndSeed)
		// bn254_optimized.Run(b, 80_000, 10, rndSeed)
		// bn254_optimized.Run(b, 100_000, 10, rndSeed)
		// bn254_optimized.Run(b, 1_000_000, 10, rndSeed)

	} else if kind == "only-bn254-crk" {

		bn254_crk.Run(b, 5_000, 10, rndSeed)
		bn254_crk.Run(b, 10_000, 10, rndSeed)
		bn254_crk.Run(b, 20_000, 10, rndSeed)
		bn254_crk.Run(b, 40_000, 10, rndSeed)
		bn254_crk.Run(b, 80_000, 10, rndSeed)
		bn254_crk.Run(b, 1_000_000, 10, rndSeed)

	} else if kind == "only-bn254-v2-wea" {

		bn254_v2_wea.Run(b, 5_000, 10, rndSeed)
		bn254_v2_wea.Run(b, 10_000, 10, rndSeed)
		// bn254_v2_wea.Run(b, 20_000, 10, rndSeed)
		// bn254_v2_wea.Run(b, 40_000, 10, rndSeed)
		// bn254_v2_wea.Run(b, 80_000, 10, rndSeed)
		// bn254_v2_wea.Run(b, 1_000_000, 10, rndSeed)

	} else if kind == "all-curves" {

		_Benchmark_Curves(b, 5_000, 10, rndSeed)
		_Benchmark_Curves(b, 10_000, 10, rndSeed)
		_Benchmark_Curves(b, 20_000, 10, rndSeed)
		_Benchmark_Curves(b, 40_000, 10, rndSeed)
		_Benchmark_Curves(b, 80_000, 10, rndSeed)
		_Benchmark_Curves(b, 100_000, 10, rndSeed)

	} else if kind == "all-results-from-paper" {
		//note: benchmark results used in the official ECPDKSAP paper

		nRandomSeeds := 10
		nRepetitions := 1

		//--------------- All curves comparison on 80k sample

		allResults := map[string]time.Duration{}
		for i := 0; i < 0; i++ {

			rndSeed := 3327 + i

			sampleSize := 80_000

			var tmp map[string]time.Duration

			tmp = bls12_377.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bls12_377.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			tmp = bls12_381.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bls12_381.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			tmp = bls24_315.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bls24_315.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			tmp = bn254.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bn254.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			tmp = bw6_633.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bw6_633.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			tmp = bw6_761.Run(b, sampleSize, nRepetitions, true, rndSeed)
			allResults["bw6_761.v2.v0-2bytes"] += tmp["v2.v0-2bytes"]

			fmt.Println("--------- Running avg. All curves comparison on 80k, nRandomSeeds:", nRandomSeeds)
			fmt.Println(allResults)
			fmt.Println()
		}

		fmt.Println("--------- Done. All curves comparison on 80k, nRandomSeeds:", nRandomSeeds)
		fmt.Println(allResults)
		fmt.Println()

		//--------------- Only BN254 optimized

		SumResults := func(prefix string, tmp map[string]time.Duration) {
			protocolVersions := []string{
				"v0.none", "v0.v0-1byte", "v0.v0-2bytes", "v0.v1-1byte", "v0.v0-11nibbles",
				"v1.none", "v1.v0-1byte", "v1.v0-2bytes", "v1.v1-1byte", "v1.v0-11nibbles",
				"v2.none", "v2.v0-1byte", "v2.v0-2bytes", "v2.v1-1byte", "v2.v0-11nibbles",
			}

			for _, pVersion := range protocolVersions {
				allResults[prefix+"."+pVersion] += tmp[pVersion]
			}
		}

		for i := 0; i < 0; i++ {

			rndSeed := 3327 + i

			// SumResults("bn254_optimized.5k", bn254_optimized.Run(b, 5_000, nRepetitions, rndSeed))
			// SumResults("bn254_optimized.10k", bn254_optimized.Run(b, 10_000, nRepetitions, rndSeed))
			// SumResults("bn254_optimized.20k", bn254_optimized.Run(b, 20_000, nRepetitions, rndSeed))
			// SumResults("bn254_optimized.40k", bn254_optimized.Run(b, 40_000, nRepetitions, rndSeed))
			// SumResults("bn254_optimized.80k", bn254_optimized.Run(b, 80_000, nRepetitions, rndSeed))
			SumResults("bn254_optimized.1mil", bn254_optimized.Run(b, 1_000_000, nRepetitions, rndSeed))

			fmt.Println("--------- Running avg. bn254_optimized: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
			fmt.Println(allResults)
			fmt.Println()
		}

		fmt.Println("--------- Done. bn254_optimized: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
		fmt.Println(allResults)
		fmt.Println()

		//--------------- Only BN254 (optimized): constant private keys (2/3)

		for i := 0; i < 0; i++ {

			rndSeed := 3327 + i

			SumResults("bn254_crk.5k", bn254_crk.Run(b, 5_000, nRepetitions, rndSeed))
			SumResults("bn254_crk.10k", bn254_crk.Run(b, 10_000, nRepetitions, rndSeed))
			SumResults("bn254_crk.20k", bn254_crk.Run(b, 20_000, nRepetitions, rndSeed))
			SumResults("bn254_crk.40k", bn254_crk.Run(b, 40_000, nRepetitions, rndSeed))
			SumResults("bn254_crk.80k", bn254_crk.Run(b, 80_000, nRepetitions, rndSeed))
			SumResults("bn254_crk.1mil", bn254_crk.Run(b, 1_000_000, nRepetitions, rndSeed))

			fmt.Println("--------- Running avg. bn254_crk: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
			fmt.Println(allResults)
			fmt.Println()
		}

		fmt.Println("--------- Done. bn254_crk: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
		fmt.Println(allResults)
		fmt.Println()

		//--------------- Only BN254 (v2 without eth address)

		for i := 0; i < nRandomSeeds; i++ {

			rndSeed := 3327 + i

			// SumResults("bn254_v2_wea.5k", bn254_v2_wea.Run(b, 5_000, nRepetitions, rndSeed))
			// SumResults("bn254_v2_wea.10k", bn254_v2_wea.Run(b, 10_000, nRepetitions, rndSeed))
			// SumResults("bn254_v2_wea.20k", bn254_v2_wea.Run(b, 20_000, nRepetitions, rndSeed))
			// SumResults("bn254_v2_wea.40k", bn254_v2_wea.Run(b, 40_000, nRepetitions, rndSeed))
			// SumResults("bn254_v2_wea.80k", bn254_v2_wea.Run(b, 80_000, nRepetitions, rndSeed))
			SumResults("bn254_v2_wea.1mil", bn254_v2_wea.Run(b, 1_000_000, nRepetitions, rndSeed))

			fmt.Println("--------- Running avg. bn254_v2_wea: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
			fmt.Println(allResults)
			fmt.Println()
		}

		fmt.Println("--------- Done. bn254_v2_wea: 5k - 1mil, nRandomSeeds:", nRandomSeeds)
		fmt.Println(allResults)
		fmt.Println()

		//--------------- Average time cost per operation

		for i := 0; i < 0; i++ {

			rndSeed := 3327 + i

			rndGen := rand.New(rand.NewSource(int64(rndSeed)))

			for j := 0; j < 1_000; j++ {
				_, v_asBigInt := _RandomPrivateKey(rndGen)
				_, _, Rj, _ := _EC_GenerateG1KeyPair(rndGen)
				_, _, _, K := _EC_GenerateG2KeyPair(rndGen)

				//-------- Shared secret calculation (recipient's side)
				var vR EC.G1Jac

				b.ResetTimer()
				vR.ScalarMultiplication(&Rj, &v_asBigInt)
				// fmt.Println("ScalarMul: ", b.Elapsed())
				allResults["ScalarMul"] += b.Elapsed()

				b.ResetTimer()
				EC.PrecomputationForFixedScalarMultiplication(&v_asBigInt)
				// fmt.Println("PrecomputationForFixedScalarMultiplication: ", b.Elapsed())
				allResults["PFSM"] += b.Elapsed()

				var table [15]EC.G1Jac

				neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix := EC.PrecomputationForFixedScalarMultiplication(&v_asBigInt)

				b.ResetTimer()
				vR.FixedScalarMultiplication(&Rj, &table, neg, k1, k2, tableElementNeeded, hiWordIndex, useMatrix)
				// fmt.Println("FixedScalarMultiplication: ", b.Elapsed())
				allResults["FSM"] += b.Elapsed()

				//-------- Pairing calculation
				b.ResetTimer()
				precomputedQLines := [][2][66]EC.LineEvaluationAff{EC.PrecomputeLines(K)}
				// fmt.Println("Precomputation for Pairing: ", b.Elapsed())
				allResults["PPair"] += b.Elapsed()

				vR_asAff := new(EC.G1Affine)

				b.ResetTimer()
				EC.PairFixedQ([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, precomputedQLines)
				// fmt.Println("Pairing: ", b.Elapsed())
				allResults["Pair"] += b.Elapsed()

				hasher := sha256.New()

				b.ResetTimer()
				hasher.Reset()
				compressed := vR_asAff.FromJacobian(&vR).Bytes()
				hash := hasher.Sum(compressed[:])
				allResults["SharedSecretHash"] += b.Elapsed()

				b.ResetTimer()
				g1, _, _, _ := EC.Generators()
				tmp := new(big.Int)
				tmp.SetBytes(hash)
				g1.ScalarMultiplicationBase(tmp)
				allResults["Protocol1.specificCalc"] += b.Elapsed()

				b.ResetTimer()
				EC.Pair([]EC.G1Affine{*vR_asAff.FromJacobian(&vR)}, []EC.G2Affine{K})
				// fmt.Println("BasicPair: ", b.Elapsed())
				allResults["BasicPair"] += b.Elapsed()

				b.ResetTimer()
				vR_asAff.FromJacobian(&vR)
				// fmt.Println("FromJacobian: ", b.Elapsed())
				allResults["FromJacobian"] += b.Elapsed()

				b.ResetTimer()
				vR_asAff.FromJacobianCoordX(&vR)
				// fmt.Println("FromJacobianCoordX: ", b.Elapsed())
				allResults["FromJacobianCoordX"] += b.Elapsed()

			}
		}
		fmt.Println("--------- Done. All done.", nRandomSeeds)
		fmt.Println(allResults)
		fmt.Println()
	}
}

func _Benchmark_Curves(b *testing.B, sampleSize int, nRepetitions int, rndSeed int) {

	bls12_377.Run(b, sampleSize, nRepetitions, true, rndSeed)
	bls12_381.Run(b, sampleSize, nRepetitions, true, rndSeed)
	bls24_315.Run(b, sampleSize, nRepetitions, true, rndSeed)
	bn254.Run(b, sampleSize, nRepetitions, true, rndSeed)
	bw6_633.Run(b, sampleSize, nRepetitions, true, rndSeed)
	bw6_761.Run(b, sampleSize, nRepetitions, true, rndSeed)
}

func _EC_GenerateG1KeyPair(r *rand.Rand) (privKey EC_fr.Element, privKey_asBigInt big.Int, pubKey EC.G1Jac, pubKeyAff EC.G1Affine) {
	g1, _, _, _ := EC.Generators()

	privKey, privKey_asBigInt = _RandomPrivateKey(r)

	pubKey.ScalarMultiplication(&g1, &privKey_asBigInt)
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
	pubKeyAff.FromJacobian(&pubKey)

	return
}

<<<<<<< HEAD
=======
func _EC_GenerateG2KeyPair(r *rand.Rand) (privKey EC_fr.Element, privKey_asBigInt big.Int, pubKey EC.G2Jac, pubKeyAff EC.G2Affine) {
	_, g2, _, _ := EC.Generators()

	privKey, privKey_asBigInt = _RandomPrivateKey(r)

	privKey.BigInt(&privKey_asBigInt)
	pubKey.ScalarMultiplication(&g2, &privKey_asBigInt)
	pubKeyAff.FromJacobian(&pubKey)

	return
}

func _RandomPrivateKey(r *rand.Rand) (privKey EC_fr.Element, privKey_asBigInt big.Int) {
	randBigInt := big.NewInt(r.Int63())
	randBigInt.Mul(randBigInt, randBigInt).Mul(randBigInt, randBigInt)
	privKey.SetBigInt(randBigInt)

	privKey.BigInt(&privKey_asBigInt)

	return
}
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28

type _CombinedMeta struct {
	Rj                *EC.G1Jac
	Rj_asAffArr       []EC.G1Affine
	ViewTagTwoBytes   uint16
	ViewTagSingleByte uint8
}
<<<<<<< HEAD


func RunAll () {

	b := new (testing.B)
	b.StartTimer()
	
	_Benchmark_tables_BN254(b)
	
}

func _Benchmark_tables_BN254(b *testing.B) {
	bn254.Run(b, 5_000, 10, true)
	bn254.Run(b, 10_000, 10, true)
	bn254.Run(b, 20_000, 10, true)
	bn254.Run(b, 40_000, 10, true)
	bn254.Run(b, 80_000, 10, true)
	bn254.Run(b, 100_000, 10, true)
	bn254.Run(b, 500_000, 10, true)
	bn254.Run(b, 1_000_000, 10, true)
	bn254.Run(b, 5_000_000, 10, true)
}

func _Benchmark_Curves(b *testing.B, sampleSize int, nRepetitions int) {

	bls12_377.Run(b, sampleSize, nRepetitions, true)
	bls12_381.Run(b, sampleSize, nRepetitions, true)
	bls24_315.Run(b, sampleSize, nRepetitions, true)
	bn254.Run(b, sampleSize, nRepetitions, true)
	bw6_633.Run(b, sampleSize, nRepetitions, true)
	bw6_761.Run(b, sampleSize, nRepetitions, true)
}
=======
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
