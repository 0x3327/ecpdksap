package main

func main() {

	const N = 5000

	// fmt.Println("\nBenchmark running with N =", N)
	// var duration time.Duration

	// fmt.Println()
	// duration = V0(N)
	// fmt.Println("V0 :: duration:", duration)
	// duration = V0_withViewTag(N)
	// fmt.Println("V0 with view tag:: duration:", duration)

	// duration = V1(N)
	// fmt.Println("V1 :: duration:", duration)

	// fmt.Println("")
}

// func V0(N int) (duration time.Duration) {

// 	K, v, Rs, _ := generalSetup(N)

// 	startTime := time.Now()

// 	for _, Rj := range Rs {
// 		utils.RecipientComputesStealthPubKey(&K, &Rj, &v)
// 	}

// 	duration = time.Since(startTime)

// 	return duration
// }

// func V0_withViewTag(N int) (duration time.Duration) {

// 	K, v, Rs, vTags := generalSetup(N)

// 	startTime := time.Now()

// 	for idx, Rj := range Rs {

// 		currVTag := utils.CalculateViewTag(&v, &Rj)

// 		if vTags[idx] == currVTag {
// 			utils.RecipientComputesStealthPubKey(&K, &Rj, &v)
// 		}
// 	}

// 	duration = time.Since(startTime)

// 	return duration
// }

// func V1(N int) (duration time.Duration) {

// 	K, v, Rs, _ := generalSetup(N)

// 	startTime := time.Now()

// 	for _, Rj := range Rs {
// 		newR := Rj
// 		utils.RecipientComputesStealthPubKey(&K, &newR, &v)
// 	}

// 	duration = time.Since(startTime)

// 	return duration
// }

// func generalSetup(N int) (K bn254.G2Affine, v fr.Element, Rs []bn254.G1Affine, vTags []string) {

// 	_, K, _ = utils.BN254_GenG2KeyPair()
// 	v, _, _ = utils.BN254_GenG1KeyPair()

// 	RsString, vTags := utils.GenRandomRs(N)

// 	var tempR bn254.G1Affine
// 	for i := 0; i < N; i++ {
// 		bytes, _ := hex.DecodeString(RsString[i])
// 		tempR.Unmarshal(bytes)
// 		Rs = append(Rs, tempR)
// 	}

// 	return K, v, Rs, vTags
// }
