package utils

import "encoding/json"

func UnmarshalUnstructredData(data json.RawMessage) (map[string]any, error) {

	var result map[string]any
	err := json.Unmarshal(data, &result)

	if err != nil {
		return nil, err
	}

	return result, nil

}
