package main

import (
<<<<<<< HEAD
<<<<<<< HEAD:go/main.go
	data_generation "ecpdksap-go/data_generation"
=======
	"ecpdksap-go/benchmark"
	"ecpdksap-go/gen_example"
>>>>>>> 1883ed2f70b6a4a34f0d24c5edd4078c2756a464:impl/main.go
=======
	"ecpdksap-go/benchmark"
	"ecpdksap-go/data_generation"
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"fmt"
	"os"
<<<<<<< HEAD
	"strings"
	"syscall/js"
)


=======
	"strconv"
	//wasm.build:::"strings"
	//wasm.build:::"syscall/js"
)

>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
func main() {

	if len(os.Args) == 1 {
		panic(`No subcommand passed - 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed!`)
	}

	subcmd := os.Args[1]

<<<<<<< HEAD
	c := make(chan string, 0)

	switch subcmd {

	case "gen-send-info":
		if len(os.Args) != 2 {
			panic(`Subcommand 'gen-send-info' takes no input params!`)
		}
		
		senderMeta := data_generation.GenerateSenderInfo("v2", "") 
		js.Global().Set("senderMeta", senderMeta)


	case "gen-recipient-info":
		if len(os.Args) != 2 {
			panic(`Subcommand 'gen-recipient-info' takes no input params!`)
		}
		
		recipientMeta := data_generation.GenerateRecipientInfo("v2", "")
		js.Global().Set("recipientMeta", recipientMeta)

=======
	switch subcmd {

>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
	case "send":
		if len(os.Args) != 3 {
			panic(`Subcommand 'send' receives all info. as one JSON input string!`)
		}
		_, _, viewTag, pubKey, addr := sender.Send(os.Args[2])
<<<<<<< HEAD
		js.Global().Set("StealthPubKey", pubKey)
		js.Global().Set("StealthAddress", addr)
		js.Global().Set("StealthViewTag", viewTag)
		
=======
		fmt.Println("Generated::: viewTag:", viewTag, "pubKey:", pubKey, "addr:", addr)
		//wasm.build:::js.Global().Set("StealthPubKey", pubKey)
		//wasm.build:::js.Global().Set("StealthAddress", addr)
		//wasm.build:::js.Global().Set("StealthViewTag", viewTag)

>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
	case "receive-scan":
		if len(os.Args) != 3 {
			panic(`Subcommand 'receive-scan' receives all info. as one JSON input string!`)
		}
		_, addrs, privKeys := recipient.Scan(os.Args[2])
<<<<<<< HEAD
		js.Global().Set("DiscoveredStealthAddrs", strings.Join(addrs, "."))
		js.Global().Set("DiscoveredStealthPrivKeys", strings.Join(privKeys, "."))
=======
		fmt.Println("Potential::: addrs:", addrs, "privateKeys", privKeys)
		//wasm.build:::js.Global().Set("DiscoveredStealthAddrs", strings.Join(addrs, "."))
		//wasm.build:::js.Global().Set("DiscoveredStealthPrivKeys", strings.Join(privKeys, "."))
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28

	case "gen-example":
		if len(os.Args) != 5 {
			panic(`Subcommand 'gen-example' needs: <version: v0 | v2> <view-tag-version: none | v0-1byte | v0-2bytes | v1-1byte> <sample-size: uint>!`)
		}
		data_generation.GenerateExample(os.Args[2], os.Args[3], os.Args[4])

<<<<<<< HEAD
	case "bench":
		if len(os.Args) != 2 {
			panic(`Subcommand 'bench' takes no arguments!`)
		}
		benchmark.RunAll()
=======
	case "gen-send-info":
		if len(os.Args) != 2 {
			panic(`Subcommand 'gen-send-info' takes no input params!`)
		}

		senderMeta := data_generation.GenerateSenderInfo() 
		fmt.Println("senderMeta", senderMeta)
		//wasm.build:::js.Global().Set("senderMeta", senderMeta)

	case "gen-recipient-info":
		if len(os.Args) != 3 {
			panic(`Subcommand 'gen-recipient-info' takes protocol version param < v0 | v1 | v2 >!`)
		}

		recipientMeta := data_generation.GenerateRecipientInfo(os.Args[2])
		fmt.Println("recipientMeta", recipientMeta)
		//wasm.build:::js.Global().Set("recipientMeta", recipientMeta)

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
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28

	default:
		fmt.Printf("\nERR: Only: 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed.\n\n")
		return
	}
<<<<<<< HEAD

	<-c
=======
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
}
