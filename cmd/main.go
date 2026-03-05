package main

import (
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/validation"
)

func printHelp() {
	fmt.Println("DIP CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  dip <command> [arguments]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  validate    Validate a decision record against a schema")
	fmt.Println("  sign        Sign a decision record (not implemented yet)")
	fmt.Println("  verify      Verify a signed decision record (not implemented yet)")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  dip validate record.json schema.json")
}

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {

	case "validate":

		if len(os.Args) < 4 {
			fmt.Println("Usage: dip validate <record.json> <schema.json>")
			return
		}

		recordPath := os.Args[2]
		schemaPath := os.Args[3]

		err := validation.Validate(recordPath, schemaPath)

		if err != nil {
			fmt.Println("Validation error:", err)
			os.Exit(1)
		}

	case "sign":
		fmt.Println("Sign command not implemented yet")

	case "verify":
		fmt.Println("Verify command not implemented yet")

	case "help":
		printHelp()

	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("")
		printHelp()
	}
}