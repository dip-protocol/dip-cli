package signing

import (
	"crypto/ed25519"
	"crypto/rand"
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
	Signature  *Signature             `json:"signature,omitempty"`
}

func Sign(recordPath string) error {

	data, err := os.ReadFile(recordPath)
	if err != nil {
		return err
	}

	var record DecisionRecord
	err = json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	// Remove signature before signing
	record.Signature = nil

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	sig := ed25519.Sign(privateKey, payload)

	record.Signature = &Signature{
		Algorithm: "ed25519",
		Value:     base64.StdEncoding.EncodeToString(sig),
	}

	signedData, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(recordPath, signedData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Record signed successfully")
	fmt.Println("Public key:", base64.StdEncoding.EncodeToString(publicKey))

	return nil
}