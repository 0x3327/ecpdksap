package benchmark

import (
	"encoding/json"
	"fmt"

	"ecpdksap-go/gen_example"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"testing"
)

func Benchmark_ThroughCLI(b *testing.B) {

	sampleSize := "1000"

	protocolVersions := []string{"v0", "v1", "v2"}
	viewTagVersions := []string{"none", "v0-1byte", "v0-2bytes", "v1-1byte"}

	for _, pVersion := range protocolVersions {

		for _, vtVersion := range viewTagVersions {

			fmt.Println("")

			sendParams, recipientParams := gen_example.GenerateExample(pVersion, vtVersion, sampleSize)

			jsonBytes, _ := json.MarshalIndent(sendParams, "", " ")
			sender.Send(string(jsonBytes))

			jsonBytes, _ = json.MarshalIndent(recipientParams, "", " ")
			recipient.Scan(string(jsonBytes))

			fmt.Println("")
		}
	}
}
