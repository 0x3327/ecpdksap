package main

import (
	"ecpdksap-bn254/gen_example"
	"ecpdksap-bn254/recipient"
	"ecpdksap-bn254/sender"
	"ecpdksap-bn254/utils"
	ecpdksap_v2 "ecpdksap-bn254/versions/v2"
	"fmt"
	"os"
)

func DBG() {
	privK, pubK := utils.GenSECP256k1G1KeyPair()

	addr := ecpdksap_v2.ComputeEthAddress(&pubK)

	fmt.Println("privK:", "0x"+privK.Text(16), "\n")
	fmt.Println("pubK:", pubK.X.Text(16)+"."+pubK.Y.Text(16))
	fmt.Println("addr:", addr)
}

func main() {

	subcmd := os.Args[1]

	switch subcmd {

	case "send":
		if len(os.Args) != 3 { panic("Subcommand `send` receives all info. as one JSON input string!") }
		sender.Send(os.Args[2])

	case "receive-scan":
		if len(os.Args) != 3 { panic("Subcommand `receive-scan` receives all info. as one JSON input string!") }
		recipient.Scan(os.Args[2])

	case "gen-example":
		if len(os.Args) != 4 { panic("Subcommand `gen-example` needs <version: v0 | v2> and <sample-size: uint>") }
		gen_example.GenerateExample(os.Args[2], os.Args[3])

	case "dbg":
		DBG()


	default:
		fmt.Printf("\nERR: only: `send` | `receive-scan`  subcommands allowed.\n\n")
		return
	}
}