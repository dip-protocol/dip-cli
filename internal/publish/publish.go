package publish

import (
	"encoding/json"
	"fmt"
	"os"
)

type RegistryRecord struct {
	ID   string                 `json:"decision_id"`
	Data map[string]interface{} `json:"data"`
}

func Publish(filePath string) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var record map[string]interface{}

	err = json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	rec := RegistryRecord{
		ID:   record["decision_id"].(string),
		Data: record,
	}

	file, err := os.OpenFile("../dip-registry/registry.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	_, err = file.Write(append(bytes, '\n'))

	if err != nil {
		return err
	}

	fmt.Println("Artifact published to registry")

	return nil
}