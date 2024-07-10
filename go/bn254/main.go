package main

import (
	"ecpdksap-bn254/recipient"
	"ecpdksap-bn254/sender"
	"fmt"
	"os"
)

func main () {

	args := os.Args

	if len(args) == 1 {
		fmt.Printf("\nERR: no `subcmd` passed!\n\n")
		return
	}

	subcmd := args[1];

	switch subcmd {

		case "send":
			sender.Send()
		case "receive-scan":
			recipient.Scan()
		case "receive-scan-using-vtag":
			recipient.ScanUsingViewTag()

		default:
			fmt.Printf("\nERR: only: `send` | `receive-scan` | `receive-scan-using-vtag` subcommands allowed.\n\n")
			return
	}
}

// k, K, _ := utils.GenG2KeyPair()
// v, V, _ := utils.GenG1KeyPair()
// fmt.Println("k:", k)
// fmt.Println("K:", K.Marshal())
// fmt.Println("v:", v)
// fmt.Println("V:", V.Marshal())
// str := hex.EncodeToString(K.Marshal())
// fmt.Println("K:", str)
// str = hex.EncodeToString(V.Marshal())
// fmt.Println("V:", str)
// // b, _ := hex.DecodeString(str)
// // fmt.Println(b)