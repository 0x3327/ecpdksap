package main

import (
	"ecpdksap-go/gen_example"
	"ecpdksap-go/recipient"
	"ecpdksap-go/sender"
	"encoding/json"
)

func main() {

	sampleSize := "1000000"

	versions := []string{"v0", "v2"}

	for _, version := range versions {
		sendParams, recipientParams := gen_example.GenerateExample(version, sampleSize)
		
		jsonBytes, _ := json.MarshalIndent(sendParams, "", " ")
		sender.Send(string(jsonBytes))

		recipientParams.WithViewTag = true
		jsonBytes, _ = json.MarshalIndent(recipientParams, "", " ")
		recipient.Scan(string(jsonBytes))

		recipientParams.WithViewTag = false
		jsonBytes, _ = json.MarshalIndent(recipientParams, "", " ")
		recipient.Scan(string(jsonBytes))
	}
}