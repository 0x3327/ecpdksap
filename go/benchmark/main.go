package main

import (
	"encoding/json"
	"fmt"

	"ecpdksap-go/gen_example"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
)

func main() {

	sampleSize := "5000"

	protocolVersions := []string{"v0", "v1", "v2"}
	viewTagVersions := []string{"none", "v0-1byte", "v0-2bytes", "v1-1byte"}

	for _, pVersion := range protocolVersions {

		sendParams, recipientParams := gen_example.GenerateExample(pVersion, "v0-1byte", sampleSize)

		for _, vtVersion := range viewTagVersions {

			fmt.Println("")

			jsonBytes, _ := json.MarshalIndent(sendParams, "", " ")
			sender.Send(string(jsonBytes))

			recipientParams.ViewTagVersion = vtVersion
			jsonBytes, _ = json.MarshalIndent(recipientParams, "", " ")
			recipient.Scan(string(jsonBytes))
		}
	}
}