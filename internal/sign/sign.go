package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"

	canonical "github.com/dip-protocol/dip-cli/internal/canonical"
)

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey string `json:"public_key"`
	Value     string `json:"value"`
}

func Sign(filePath string) (string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var record map[string]interface{}

	err = json.Unmarshal(data, &record)
	if err != nil {
		return "", err
	}

	// Remove existing signature before signing
	delete(record, "signature")

	// Canonicalize JSON
	canonicalBytes, err := canonical.Canonicalize(record)
	if err != nil {
		return "", err
	}

	// Generate keypair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}

	// Sign canonical bytes
	signatureBytes := ed25519.Sign(privateKey, canonicalBytes)

	sig := Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(publicKey),
		Value:     base64.StdEncoding.EncodeToString(signatureBytes),
	}

	// Attach signature object
	record["signature"] = sig

	out, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filePath, out, 0644)
	if err != nil {
		return "", err
	}

	return sig.PublicKey, nil
}