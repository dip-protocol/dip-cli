package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("DIP CLI")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  dip verify <artifact.json>")
		fmt.Println("  dip sign <artifact.json>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "verify":
		fmt.Println("Verifying artifact...")

	case "sign":
		fmt.Println("Signing artifact...")

	default:
		fmt.Println("Unknown command")
	}
}