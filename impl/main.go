package main

import (
<<<<<<< HEAD:go/main.go
	data_generation "ecpdksap-go/data_generation"
=======
	"ecpdksap-go/benchmark"
	"ecpdksap-go/gen_example"
>>>>>>> 1883ed2f70b6a4a34f0d24c5edd4078c2756a464:impl/main.go
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"fmt"
	"os"
	"strings"
	"syscall/js"
)


func main() {

	if len(os.Args) == 1 {
		panic(`No subcommand passed - 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed!`)
	}

	subcmd := os.Args[1]

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

	case "send":
		if len(os.Args) != 3 {
			panic(`Subcommand 'send' receives all info. as one JSON input string!`)
		}
		_, _, viewTag, pubKey, addr := sender.Send(os.Args[2])
		js.Global().Set("StealthPubKey", pubKey)
		js.Global().Set("StealthAddress", addr)
		js.Global().Set("StealthViewTag", viewTag)
		
	case "receive-scan":
		if len(os.Args) != 3 {
			panic(`Subcommand 'receive-scan' receives all info. as one JSON input string!`)
		}
		_, addrs, privKeys := recipient.Scan(os.Args[2])
		js.Global().Set("DiscoveredStealthAddrs", strings.Join(addrs, "."))
		js.Global().Set("DiscoveredStealthPrivKeys", strings.Join(privKeys, "."))

	case "gen-example":
		if len(os.Args) != 5 {
			panic(`Subcommand 'gen-example' needs: <version: v0 | v2> <view-tag-version: none | v0-1byte | v0-2bytes | v1-1byte> <sample-size: uint>!`)
		}
		data_generation.GenerateExample(os.Args[2], os.Args[3], os.Args[4])

	case "bench":
		if len(os.Args) != 2 {
			panic(`Subcommand 'bench' takes no arguments!`)
		}
		benchmark.RunAll()

	default:
		fmt.Printf("\nERR: Only: 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed.\n\n")
		return
	}

	<-c
}
