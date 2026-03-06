package main

import (
	"fmt"
	"os"

	publishpkg "github.com/dip-protocol/dip-cli/internal/publish"
	signpkg "github.com/dip-protocol/dip-cli/internal/sign"
	valpkg "github.com/dip-protocol/dip-cli/internal/validation"
	verifypkg "github.com/dip-protocol/dip-cli/internal/verify"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("DIP CLI")
		fmt.Println("Usage: dip <command>")
		fmt.Println("Commands: validate, sign, verify, publish")
		return
	}

	command := os.Args[1]

	switch command {

	case "validate":

		if len(os.Args) < 4 {
			fmt.Println("Usage: dip validate <record.json> <schema.json>")
			return
		}

		record := os.Args[2]
		schema := os.Args[3]

		err := valpkg.Validate(record, schema)

		if err != nil {
			fmt.Println("Validation error:", err)
			os.Exit(1)
		}

		fmt.Println("Record is valid according to schema")

	case "sign":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip sign <record.json>")
			return
		}

		record := os.Args[2]

		pub, err := signpkg.Sign(record)

		if err != nil {
			fmt.Println("Sign error:", err)
			os.Exit(1)
		}

		fmt.Println("Record signed successfully")
		fmt.Println("Public key:", pub)

	case "verify":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip verify <record.json>")
			return
		}

		record := os.Args[2]

		err := verifypkg.Verify(record)

		if err != nil {
			fmt.Println("Signature INVALID:", err)
			os.Exit(1)
		}

		fmt.Println("Signature valid")

	case "publish":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip publish <record.json>")
			return
		}

		record := os.Args[2]

		err := publishpkg.Publish(record)

		if err != nil {
			fmt.Println("Publish error:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command")
	}
}