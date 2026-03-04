package verify

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/hash"
	"github.com/dip-protocol/dip-cli/internal/schema"
)

func VerifyArtifact(path string) error {

	// ---- Step 1: Schema Validation ----

	schemaPath := "D:/Conf/dip-spec/schemas/artifact.schema.json"

	err := schema.ValidateArtifact(path, schemaPath)
	if err != nil {
		return err
	}

	// ---- Step 2: Read Artifact ----

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var artifact map[string]interface{}

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		return err
	}

	// ---- Step 3: Extract Stored Fields ----

	storedHash, ok := artifact["artifact_hash"].(string)
	if !ok {
		return fmt.Errorf("artifact_hash missing or invalid")
	}

	signatureHex, ok := artifact["signature"].(string)
	if !ok {
		return fmt.Errorf("signature missing")
	}

	publicKeyHex, ok := artifact["public_key"].(string)
	if !ok {
		return fmt.Errorf("public_key missing")
	}

	// ---- Step 4: Remove Non-Hashed Fields ----

	delete(artifact, "artifact_hash")
	delete(artifact, "signature")
	delete(artifact, "public_key")

	// ---- Step 5: Compute Canonical Hash ----

	computedHash, err := hash.ComputeCanonical(artifact)
	if err != nil {
		return err
	}

	fmt.Println("Computed Hash:", computedHash)
	fmt.Println("Artifact Hash:", storedHash)

	if computedHash != storedHash {
		return fmt.Errorf("artifact hash mismatch")
	}

	fmt.Println("Artifact integrity verified")

	// ---- Step 6: Verify Signature ----

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return err
	}

	valid := ed25519.Verify(publicKey, []byte(storedHash), signature)

	if !valid {
		return fmt.Errorf("signature verification failed")
	}

	fmt.Println("Signature verification successful")

	return nil
}