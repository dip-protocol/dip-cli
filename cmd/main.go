package main

import (
	"fmt"
	"os"

	signpkg "github.com/dip-protocol/dip-cli/internal/signing"
	verpkg "github.com/dip-protocol/dip-cli/internal/verify"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("DIP CLI")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  dip sign <record.json>")
		fmt.Println("  dip verify <record.json> <public-key>")
		return
	}

	switch os.Args[1] {

	case "sign":

		if len(os.Args) < 3 {
			fmt.Println("Usage: dip sign <record.json>")
			return
		}

		err := signpkg.Sign(os.Args[2])
		if err != nil {
			fmt.Println("Signing error:", err)
			os.Exit(1)
		}

	case "verify":

		if len(os.Args) < 4 {
			fmt.Println("Usage: dip verify <record.json> <public-key>")
			return
		}

		err := verpkg.Verify(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println("Verification error:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Unknown command")
	}
}