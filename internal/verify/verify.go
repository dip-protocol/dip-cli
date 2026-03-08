package verify

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

func Verify(filePath string) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var record map[string]interface{}

	err = json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	sigObj, ok := record["signature"].(map[string]interface{})
	if !ok {
		return errors.New("invalid signature format")
	}

	algorithm, ok := sigObj["algorithm"].(string)
	if !ok {
		return errors.New("invalid signature algorithm")
	}

	if algorithm != "ed25519" {
		return errors.New("unsupported signature algorithm")
	}

	pubKeyStr, ok := sigObj["public_key"].(string)
	if !ok {
		return errors.New("invalid public key format")
	}

	sigStr, ok := sigObj["value"].(string)
	if !ok {
		return errors.New("invalid signature format")
	}

	sig := Signature{
		Algorithm: algorithm,
		PublicKey: pubKeyStr,
		Value:     sigStr,
	}

	artifactID, ok := record["artifact_id"].(string)
	if !ok {
		return errors.New("artifact_id missing")
	}

	// Remove signature before verification
	delete(record, "signature")

	// Canonicalize artifact
	canonicalBytes, err := canonical.Canonicalize(record)
	if err != nil {
		return err
	}

	// Recompute artifact_id
	hash := sha256.Sum256(canonicalBytes)
	expectedID := hex.EncodeToString(hash[:])

	if artifactID != expectedID {
		return errors.New("artifact_id mismatch")
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(sig.PublicKey)
	if err != nil {
		return err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sig.Value)
	if err != nil {
		return err
	}

	valid := ed25519.Verify(pubKeyBytes, canonicalBytes, sigBytes)

	if !valid {
		return errors.New("signature verification failed")
	}

	return nil
}