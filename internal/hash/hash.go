package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

func ComputeCanonical(data map[string]interface{}) (string, error) {

	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	ordered := make(map[string]interface{})

	for _, k := range keys {
		ordered[k] = data[k]
	}

	bytes, err := json.Marshal(ordered)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bytes)

	return hex.EncodeToString(hash[:]), nil
}