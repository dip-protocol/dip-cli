package canonical

import (
	"encoding/json"
	"sort"
)

func Canonicalize(v interface{}) ([]byte, error) {
	return json.Marshal(sortMap(v))
}

func sortMap(v interface{}) interface{} {

	switch val := v.(type) {

	case map[string]interface{}:

		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		sorted := make(map[string]interface{})

		for _, k := range keys {
			sorted[k] = sortMap(val[k])
		}

		return sorted

	case []interface{}:

		for i := range val {
			val[i] = sortMap(val[i])
		}

		return val
	}

	return v
}