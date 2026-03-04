package sign

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/hash"
)

func SignArtifact(path string) error {

	// ---- Read Artifact ----

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var artifact map[string]interface{}

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		return err
	}

	// ---- Remove fields before hashing ----

	delete(artifact, "artifact_hash")
	delete(artifact, "signature")
	delete(artifact, "public_key")

	// ---- Compute Canonical Hash ----

	computedHash, err := hash.ComputeCanonical(artifact)
	if err != nil {
		return err
	}

	// ---- Add Hash to Artifact ----

	artifact["artifact_hash"] = computedHash

	// ---- Generate Ed25519 Key Pair ----

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// ---- Sign the Hash ----

	signature := ed25519.Sign(privateKey, []byte(computedHash))

	artifact["signature"] = hex.EncodeToString(signature)
	artifact["public_key"] = hex.EncodeToString(publicKey)

	// ---- Write Artifact Back ----

	output, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Artifact signed successfully")

	return nil
}