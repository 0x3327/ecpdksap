package main

import (
	"ecpdksap-go/benchmark"
	"ecpdksap-go/data_generation"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"fmt"
	"os"
	"strconv"
	//wasm.build:::"syscall/js"
)

func main() {

	if len(os.Args) == 1 {
		panic(`No subcommand passed - 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed!`)
	}

	subcmd := os.Args[1]

	switch subcmd {

	case "send":
		if len(os.Args) != 3 {
			panic(`Subcommand 'send' receives all info. as one JSON input string!`)
		}
		_, _, viewTag, pubKey, addr := sender.Send(os.Args[2])
		fmt.Println("Generated::: viewTag:", viewTag, "pubKey:", pubKey, "addr:", addr)
		//wasm.build:::js.Global().Set("StealthPubKey", pubKey)
		//wasm.build:::js.Global().Set("StealthAddress", addr)
		//wasm.build:::js.Global().Set("StealthViewTag", viewTag)

	case "receive-scan":
		if len(os.Args) != 3 {
			panic(`Subcommand 'receive-scan' receives all info. as one JSON input string!`)
		}
		_, addrs, privKeys := recipient.Scan(os.Args[2])
		fmt.Println("Potential::: addrs:", addrs, "privateKeys", privKeys)
		//wasm.build:::js.Global().Set("DiscoveredStealthAddrs", strings.Join(addrs, "."))
		//wasm.build:::js.Global().Set("DiscoveredStealthPrivKeys", strings.Join(privKeys, "."))

	case "gen-example":
		if len(os.Args) != 5 {
			panic(`Subcommand 'gen-example' needs: <version: v0 | v2> <view-tag-version: none | v0-1byte | v0-2bytes | v1-1byte> <sample-size: uint>!`)
		}
		data_generation.GenerateExample(os.Args[2], os.Args[3], os.Args[4])

	case "gen-send-info":
		if len(os.Args) != 2 {
			panic(`Subcommand 'gen-send-info' takes no input params!`)
		}

		senderMeta := data_generation.GenerateSenderInfo() 
		fmt.Println("senderMeta", senderMeta)

		//wasm.build::: senderMeta := data_generation.GenerateSenderInfo("v2", "") 
		//wasm.build::: js.Global().Set("senderMeta", senderMeta)

	case "gen-recipient-info":
		if len(os.Args) != 3 {
			panic(`Subcommand 'gen-recipient-info' takes protocol version param < v0 | v1 | v2 >!`)
		}

		recipientMeta := data_generation.GenerateRecipientInfo(os.Args[2])
		fmt.Println("recipientMeta", recipientMeta)
		//wasm.build::: js.Global().Set("recipientMeta", recipientMeta)

	case "bench":
		if len(os.Args) < 3 {
			panic(`Subcommand 'bench' takes one argument <only-bn254 | all-curves>!`)
		}

		if len(os.Args) == 4 {
			fmt.Println("Received str:", os.Args[3])
			seed, _ := strconv.Atoi(os.Args[3])
			fmt.Println("Converted str to int:", seed)
			benchmark.RunBench(os.Args[2], seed)
		} else {
			defaultSeed := 12318726
			benchmark.RunBench(os.Args[2], defaultSeed)
		}

	default:
		fmt.Printf("\nERR: Only: 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed.\n\n")
		return
	}
}
