package cmd

import "encoding/json"

func jsonMarshal(value any) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}
