package utils

import (
	"encoding/json"
	"fmt"
)

func PrintMapAsJSON(data map[string]interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ") // Indent with spaces
	if err != nil {
		fmt.Println("error marshalling map to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}