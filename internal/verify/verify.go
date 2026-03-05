package verify

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type Signature struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

type DecisionRecord struct {
	Version    string                 `json:"version"`
	DecisionID string                 `json:"decision_id"`
	Timestamp  string                 `json:"timestamp"`
	Inputs     map[string]interface{} `json:"inputs"`
	Outputs    map[string]interface{} `json:"outputs"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Signature  *Signature             `json:"signature"`
}

func Verify(recordPath string, publicKeyBase64 string) error {

	data, err := os.ReadFile(recordPath)
	if err != nil {
		return err
	}

	var record DecisionRecord
	err = json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	if record.Signature == nil {
		return fmt.Errorf("no signature found in record")
	}

	// Decode signature
	sigBytes, err := base64.StdEncoding.DecodeString(record.Signature.Value)
	if err != nil {
		return err
	}

	// Decode public key
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return err
	}

	// Remove signature before verification
	record.Signature = nil

	// Create deterministic payload structure
	payloadStruct := struct {
		Version    string                 `json:"version"`
		DecisionID string                 `json:"decision_id"`
		Timestamp  string                 `json:"timestamp"`
		Inputs     map[string]interface{} `json:"inputs"`
		Outputs    map[string]interface{} `json:"outputs"`
		Metadata   map[string]interface{} `json:"metadata,omitempty"`
	}{
		record.Version,
		record.DecisionID,
		record.Timestamp,
		record.Inputs,
		record.Outputs,
		record.Metadata,
	}

	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		return err
	}

	valid := ed25519.Verify(pubKeyBytes, payload, sigBytes)

	if valid {
		fmt.Println("Signature valid")
	} else {
		fmt.Println("Signature INVALID")
	}

	return nil
}