package canonical

import (
	"bytes"
	"encoding/json"
	"sort"
)

func Canonicalize(v interface{}) ([]byte, error) {

	buf := &bytes.Buffer{}

	err := writeCanonical(buf, v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeCanonical(buf *bytes.Buffer, v interface{}) error {

	switch val := v.(type) {

	case map[string]interface{}:

		keys := make([]string, 0, len(val))

		for k := range val {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		buf.WriteByte('{')

		for i, k := range keys {

			if i > 0 {
				buf.WriteByte(',')
			}

			keyBytes, err := json.Marshal(k)
			if err != nil {
				return err
			}

			buf.Write(keyBytes)
			buf.WriteByte(':')

			err = writeCanonical(buf, val[k])
			if err != nil {
				return err
			}
		}

		buf.WriteByte('}')

	case []interface{}:

		buf.WriteByte('[')

		for i, item := range val {

			if i > 0 {
				buf.WriteByte(',')
			}

			err := writeCanonical(buf, item)
			if err != nil {
				return err
			}
		}

		buf.WriteByte(']')

	default:

		b, err := json.Marshal(val)
		if err != nil {
			return err
		}

		buf.Write(b)
	}

	return nil
}