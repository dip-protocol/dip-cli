package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/registry"
	"github.com/dip-protocol/dip-cli/internal/sign"
	"github.com/dip-protocol/dip-cli/internal/verify"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  dip sign <artifact.json>")
		fmt.Println("  dip verify <artifact.json>")
		fmt.Println("  dip publish <artifact.json>")
		return
	}

	command := os.Args[1]

	switch command {

	case "sign":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip sign <artifact.json>")
			return
		}

		err := sign.SignArtifact(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}

	case "verify":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip verify <artifact.json>")
			return
		}

		err := verify.VerifyArtifact(os.Args[2])
		if err != nil {
			fmt.Println("Verification error:", err)
		}

	case "publish":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip publish <artifact.json>")
			return
		}

		path := os.Args[2]

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var artifact map[string]interface{}

		err = json.Unmarshal(data, &artifact)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		hash, ok := artifact["artifact_hash"].(string)
		if !ok {
			fmt.Println("artifact_hash missing")
			return
		}

		err = registry.PublishArtifact(hash)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

	default:

		fmt.Println("Unknown command:", command)
		fmt.Println("Available commands:")
		fmt.Println("  sign")
		fmt.Println("  verify")
		fmt.Println("  publish")
	}
}