package verify

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
)

type Signature struct {
	Algorithm string `json:"algorithm"`
	PublicKey string `json:"public_key"`
	Value     string `json:"value"`
}

func Verify(filePath string) error {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var record map[string]interface{}

	err = json.Unmarshal(data, &record)
	if err != nil {
		return err
	}

	sigObj, ok := record["signature"].(map[string]interface{})
	if !ok {
		return errors.New("invalid signature format")
	}

	sig := Signature{
		Algorithm: sigObj["algorithm"].(string),
		PublicKey: sigObj["public_key"].(string),
		Value:     sigObj["value"].(string),
	}

	delete(record, "signature")

	canonical, err := json.Marshal(record)
	if err != nil {
		return err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(sig.PublicKey)
	if err != nil {
		return err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sig.Value)
	if err != nil {
		return err
	}

	valid := ed25519.Verify(pubKeyBytes, canonical, sigBytes)

	if !valid {
		return errors.New("signature verification failed")
	}

	return nil
}