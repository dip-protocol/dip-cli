package registry

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entry struct {
	ArtifactHash string `json:"artifact_hash"`
}

func PublishArtifact(hash string) error {

	registryPath := "D:/Conf/dip-registry/log.json"

	var entries []Entry

	data, err := os.ReadFile(registryPath)
	if err == nil {
		json.Unmarshal(data, &entries)
	}

	entry := Entry{
		ArtifactHash: hash,
	}

	entries = append(entries, entry)

	output, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(registryPath, output, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Artifact published to registry")

	return nil
}