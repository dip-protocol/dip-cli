package proof

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dip-protocol/dip-cli/internal/merkle"
)

type LogEntry struct {
	Artifact string `json:"artifact"`
	Hash     string `json:"hash"`
}

type Proof struct {
	ArtifactHash string   `json:"artifact_hash"`
	ProofPath    []string `json:"proof_path"`
	Root         string   `json:"root"`
}

func GenerateProof(artifactPath string) error {

	registryDir := "D:/Conf/dip-registry"

	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return err
	}

	targetHash := merkle.Hash(data)

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

	index := -1

	for i, entry := range log {
		if entry.Hash == targetHash {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Println("Artifact not found in registry")
		return nil
	}

	var path []string

	for i, entry := range log {
		if i != index {
			path = append(path, entry.Hash)
		}
	}

	root := ""

	for _, entry := range log {

		if root == "" {
			root = entry.Hash
		} else {
			root = merkle.Combine(root, entry.Hash)
		}

	}

	proof := Proof{
		ArtifactHash: targetHash,
		ProofPath:    path,
		Root:         root,
	}

	output, _ := json.MarshalIndent(proof, "", "  ")

	proofFile := "proof.json"

	err = os.WriteFile(proofFile, output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Proof generated:", proofFile)

	return nil
}