package sign

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

	canonical "github.com/dip-protocol/dip-cli/internal/canonical"
)

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey string `json:"public_key"`
	Value     string `json:"value"`
}

func loadPrivateKey(path string) (ed25519.PrivateKey, error) {

	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(string(keyBytes))
	if err != nil {
		return nil, err
	}

	if len(key) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid ed25519 private key length")
	}

	return ed25519.PrivateKey(key), nil
}

// Sign generates a DIP signature for the artifact.
//
// Signing procedure defined by the DIP protocol:
//
// 1. Remove the signature field
// 2. Canonicalize artifact
// 3. Compute artifact_id = SHA256(canonical_bytes)
// 4. Insert artifact_id
// 5. Canonicalize again
// 6. Sign canonical bytes using Ed25519
func Sign(filePath string, keyPath string) (string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var record map[string]interface{}

	err = json.Unmarshal(data, &record)
	if err != nil {
		return "", err
	}

	// Ensure protocol version
	record["artifact_version"] = "1"

	// Remove existing signature
	delete(record, "signature")

	// First canonicalization (without artifact_id)
	canonicalBytes, err := canonical.Canonicalize(record)
	if err != nil {
		return "", err
	}

	// Compute artifact_id
	hash := sha256.Sum256(canonicalBytes)
	artifactID := hex.EncodeToString(hash[:])

	record["artifact_id"] = artifactID

	// Canonicalize again including artifact_id
	canonicalBytes, err = canonical.Canonicalize(record)
	if err != nil {
		return "", err
	}

	// Load signing key
	privateKey, err := loadPrivateKey(keyPath)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	// Sign canonical artifact
	signatureBytes := ed25519.Sign(privateKey, canonicalBytes)

	sig := Signature{
		Algorithm: "ed25519",
		PublicKey: base64.StdEncoding.EncodeToString(publicKey),
		Value:     base64.StdEncoding.EncodeToString(signatureBytes),
	}

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