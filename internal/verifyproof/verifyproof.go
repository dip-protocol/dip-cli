package verifyproof

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/merkle"
)

type Artifact struct {
	ArtifactID string `json:"artifact_id"`
}

type ProofNode struct {
	Hash     string `json:"hash"`
	Position string `json:"position"`
}

type Proof struct {
	ArtifactHash string      `json:"artifact_hash"`
	ProofPath    []ProofNode `json:"proof_path"`
	Root         string      `json:"root"`
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

	if proof.ArtifactHash != artifact.ArtifactID {
		return errors.New("artifact hash does not match proof")
	}

	hash := artifact.ArtifactID

	for _, node := range proof.ProofPath {

		if node.Position == "left" {

			hash = merkle.Combine(node.Hash, hash)

		} else if node.Position == "right" {

			hash = merkle.Combine(hash, node.Hash)

		} else {

			return errors.New("invalid proof position")
		}
	}

	if hash != proof.Root {

		fmt.Println("Proof verification: INVALID")
		return errors.New("merkle root mismatch")
	}

	fmt.Println("Proof verification: VALID")

	return nil
}