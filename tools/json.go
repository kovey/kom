package tools

import "encoding/json"

func ArrToJson[T any](data []T) string {
	buf, err := json.Marshal(data)
	if err != nil {
		return "[]"
	}

	return string(buf)
}

func MapToJson[T any](data map[string]T) string {
	buf, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(buf)
}

func StructToJson[T any](data T) string {
	buf, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(buf)
}
