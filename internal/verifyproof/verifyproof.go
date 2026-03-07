package verifyproof

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/merkle"
)

type Artifact struct {
	ArtifactID string `json:"artifact_id"`
}

type Proof struct {
	ArtifactHash string   `json:"artifact_hash"`
	ProofPath    []string `json:"proof_path"`
	Root         string   `json:"root"`
}

func VerifyProof(artifactPath string, proofPath string) error {

	artifactData, err := os.ReadFile(artifactPath)
	if err != nil {
		return err
	}

	var artifact Artifact

	err = json.Unmarshal(artifactData, &artifact)
	if err != nil {
		return err
	}

	proofData, err := os.ReadFile(proofPath)
	if err != nil {
		return err
	}

	var proof Proof

	err = json.Unmarshal(proofData, &proof)
	if err != nil {
		return err
	}

	hash := artifact.ArtifactID

	for _, sibling := range proof.ProofPath {

		hash = merkle.Combine(hash, sibling)

	}

	if hash == proof.Root {

		fmt.Println("Proof verification: VALID")

	} else {

		fmt.Println("Proof verification: INVALID")

	}

	return nil
}