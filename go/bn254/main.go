package main

import (
	"ecpdksap-bn254/recipient"
	"ecpdksap-bn254/sender"
	"fmt"
	"os"
)

func MainSend (jsonInputString string) {

	fmt.Println(jsonInputString)
}


func main () {

	args := os.Args

	if len(args) == 1 {
		fmt.Printf("\nERR: no args passed!\n\n")
		return
	}

	subcmd := args[1];

	fmt.Println(args[2])

	if subcmd != "send" && subcmd != "receive-scan" {
		fmt.Printf("\nERR: only 'send' and 'receive-scan' subcommands allowed.\n\n")
		return
	}

	if subcmd == "send" {
		sender.Send()
	}

	if subcmd == "receive-scan" {
		recipient.Scan()
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