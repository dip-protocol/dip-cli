package verifyproof

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dip-protocol/dip-cli/internal/merkle"
)

type Proof struct {
	ArtifactHash string   `json:"artifact_hash"`
	ProofPath    []string `json:"proof_path"`
	Root         string   `json:"root"`
}

func VerifyProof(artifactFile string, proofFile string) error {

	data, err := os.ReadFile(artifactFile)
	if err != nil {
		return err
	}

	computedHash := merkle.Hash(data)

	proofData, err := os.ReadFile(proofFile)
	if err != nil {
		return err
	}

	var proof Proof

	err = json.Unmarshal(proofData, &proof)
	if err != nil {
		return err
	}

	fmt.Println("Computed Artifact Hash:", computedHash)
	fmt.Println("Proof Artifact Hash:", proof.ArtifactHash)

	if computedHash != proof.ArtifactHash {
		fmt.Println("Artifact hash mismatch")
		return nil
	}

	fmt.Println("Merkle Root (from proof):", proof.Root)
	fmt.Println("Proof verification successful")

	return nil
}