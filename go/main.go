package main

import (
	"ecpdksap-go/gen_example"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"ecpdksap-go/utils"
	ecpdksap_v2 "ecpdksap-go/versions/v2"
	"fmt"
	"os"
)

func DBG() {
	privK, pubK := utils.SECP256k_Gen1G1KeyPair()

	addr := ecpdksap_v2.ComputeEthAddress(&pubK)

	fmt.Println("privK:", "0x"+privK.Text(16))
	fmt.Println("pubK:", pubK.X.Text(16)+"."+pubK.Y.Text(16))
	fmt.Println("addr:", addr)
}

func main() {

	subcmd := os.Args[1]

	switch subcmd {

	case "send":
		if len(os.Args) != 3 { panic(`Subcommand 'send' receives all info. as one JSON input string!`) }
		sender.Send(os.Args[2])

	case "receive-scan":
		if len(os.Args) != 3 { panic(`Subcommand 'receive-scan' receives all info. as one JSON input string!`) }
		recipient.Scan(os.Args[2])

	case "gen-example":
		if len(os.Args) != 5 { 
			panic(`Subcommand 'gen-example' needs: <version: v0 | v2> <view-tag-version: none | v0-1byte | v0-2bytes | v1-1byte> <sample-size: uint>!`) 
		}
		gen_example.GenerateExample(os.Args[2], os.Args[3], os.Args[4])

	case "dbg":
		DBG()


	default:
		fmt.Printf("\nERR: Only: 'send' | 'receive-scan' | 'gen-example' subcommands allowed.\n\n")
		return
	}
}