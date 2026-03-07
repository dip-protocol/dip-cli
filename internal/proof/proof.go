package proof

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/dip-protocol/dip-cli/internal/merkle"
)

type Artifact struct {
	ArtifactID string `json:"artifact_id"`
}

type LogEntry struct {
	ArtifactHash string `json:"artifact_hash"`
}

type Proof struct {
	ArtifactHash string   `json:"artifact_hash"`
	ProofPath    []string `json:"proof_path"`
	Root         string   `json:"root"`
}

func GenerateProof(artifactPath string, registryDir string) error {

	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return err
	}

	var artifact Artifact

	err = json.Unmarshal(data, &artifact)
	if err != nil {
		return err
	}

	targetHash := artifact.ArtifactID

	logPath := filepath.Join(registryDir, "log.json")

	logData, err := os.ReadFile(logPath)
	if err != nil {
		return err
	}

	var log []LogEntry

	err = json.Unmarshal(logData, &log)
	if err != nil {
		return err
	}

	var leaves []string

	for _, entry := range log {
		leaves = append(leaves, entry.ArtifactHash)
	}

	// Deterministic ordering
	sort.Strings(leaves)

	index := -1

	for i, hash := range leaves {
		if hash == targetHash {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Artifact not found in registry")
		return nil
	}

	var proofPath []string

	level := leaves
	pos := index

	for len(level) > 1 {

		var nextLevel []string

		for i := 0; i < len(level); i += 2 {

			left := level[i]

			right := left
			if i+1 < len(level) {
				right = level[i+1]
			}

			parent := merkle.Combine(left, right)
			nextLevel = append(nextLevel, parent)

			if i == pos || i+1 == pos {

				if pos == i {
					proofPath = append(proofPath, right)
				} else {
					proofPath = append(proofPath, left)
				}

				pos = len(nextLevel) - 1
			}
		}

		level = nextLevel
	}

	root := level[0]

	proof := Proof{
		ArtifactHash: targetHash,
		ProofPath:    proofPath,
		Root:         root,
	}

	output, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("proof.json", output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Proof generated: proof.json")

	return nil
}