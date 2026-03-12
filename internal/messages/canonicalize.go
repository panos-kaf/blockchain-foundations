package messages

import (
	"bytes"
	"encoding/json"
	"sort"
)

// Canonicalize takes an arbitrary Go value, 
// marshals it to JSON, 
// and then re-marshals it in a canonical form 
// with sorted keys and consistent formatting. 
func canonicalizeJSON(v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := bytes.NewBufferString("{")
		for i, k := range keys {
			if i > 0 {
				buf.WriteString(",")
			}
			keyBytes, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			buf.Write(keyBytes)
			buf.WriteString(":")
			valueBytes, err := canonicalizeJSON(val[k])
			if err != nil {
				return nil, err
			}
			buf.Write(valueBytes)

		}
		buf.WriteString("}")
		return buf.Bytes(), nil

	case []interface{}:
		buf := bytes.NewBufferString("[")
		for i, elem := range val {
			if i > 0 {
				buf.WriteString(",")
			}
			elemBytes, err := canonicalizeJSON(elem)
			if err != nil {
				return nil, err
			}
			buf.Write(elemBytes)
		}
		buf.WriteString("]")
		return buf.Bytes(), nil

	default:
		return json.Marshal(val)
	}

}

func Canonicalize(v interface{}) (string, error) {
	var obj interface{}
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		return "", err
	}
	canon, err := canonicalizeJSON(obj)
	if err != nil {
		return "", err
	}
	return string(canon), nil
}