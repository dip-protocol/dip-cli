package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// ComputeArtifactID calculates the deterministic artifact identifier.
//
// According to the DIP protocol specification:
//
// artifact_id = SHA256(canonical_artifact_bytes)
//
// The input must be the canonical JSON representation of the artifact
// (with the signature field removed).
func ComputeArtifactID(data []byte) string {

	digest := sha256.Sum256(data)

	return hex.EncodeToString(digest[:])
}