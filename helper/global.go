package helper

import "encoding/json"

func JsonToStruct(data []byte, result any) error {
	return json.Unmarshal(data, &result)
}
