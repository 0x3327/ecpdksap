package main

import (
	"ecpdksap-go/benchmark"
	"ecpdksap-go/gen_example"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"fmt"
	"os"
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
		sender.Send(os.Args[2])

	case "receive-scan":
		if len(os.Args) != 3 {
			panic(`Subcommand 'receive-scan' receives all info. as one JSON input string!`)
		}
		recipient.Scan(os.Args[2])

	case "gen-example":
		if len(os.Args) != 5 {
			panic(`Subcommand 'gen-example' needs: <version: v0 | v2> <view-tag-version: none | v0-1byte | v0-2bytes | v1-1byte> <sample-size: uint>!`)
		}
		gen_example.GenerateExample(os.Args[2], os.Args[3], os.Args[4])

	case "bench":
		if len(os.Args) != 2 {
			panic(`Subcommand 'bench' takes no arguments!`)
		}
		benchmark.RunAll()

	default:
		fmt.Printf("\nERR: Only: 'send' | 'receive-scan' | 'gen-example' | 'bench' subcommands allowed.\n\n")
		return
	}
}
